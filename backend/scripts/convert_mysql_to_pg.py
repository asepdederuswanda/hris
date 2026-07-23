#!/usr/bin/env python3
"""
Convert MySQL DDL migration files to PostgreSQL syntax.
Reads from mysql/ directory, writes to postgres/ directory.
"""
import os
import re
import shutil

MIGRATIONS_DIR = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "internal/pkg/migrator/migrations/tenant"
)

MYSQL_DIR = os.path.join(MIGRATIONS_DIR, "mysql")
PG_DIR = os.path.join(MIGRATIONS_DIR, "postgres")

def convert_sql(content):
    """Convert MySQL DDL to PostgreSQL-compatible DDL."""
    
    # 1. Remove ENGINE clause
    content = re.sub(
        r'\s*ENGINE\s*=\s*\w+(\s+DEFAULT\s+CHARSET\s*=\s*\w+(\s+COLLATE\s*=\s*\w+)?)?',
        '',
        content
    )
    
    # 2. Replace TINYINT(1) with SMALLINT
    content = re.sub(r'\bTINYINT\s*\(\s*1\s*\)', 'SMALLINT', content)
    
    # 3. Replace standalone TINYINT with SMALLINT
    content = re.sub(r'\bTINYINT\b(?!\s*\()', 'SMALLINT', content)
    
    # 4. Remove all UNSIGNED keywords (PostgreSQL doesn't support UNSIGNED)
    # Handles INT UNSIGNED, BIGINT UNSIGNED, SMALLINT UNSIGNED, etc.
    content = re.sub(r'\b(INT|BIGINT|SMALLINT|TINYINT)\s+UNSIGNED\b', r'\1', content, flags=re.IGNORECASE)
    
    # 5. Replace ENUM with VARCHAR(255)
    # Must handle nested parentheses in enum values e.g. 'Dampak (Hasil)'
    # Using manual parsing instead of regex for balanced parens
    def replace_enums(text):
        result = []
        i = 0
        while i < len(text):
            # Find ENUM( pattern (case-insensitive)
            if text[i:i+5].upper() == 'ENUM(':
                start = i
                i += 5  # skip 'ENUM('
                paren_depth = 1
                while i < len(text) and paren_depth > 0:
                    if text[i] == '(':
                        paren_depth += 1
                    elif text[i] == ')':
                        paren_depth -= 1
                    i += 1
                # Replace the whole ENUM(...) with VARCHAR(255)
                result.append('VARCHAR(255)')
            else:
                result.append(text[i])
                i += 1
        return ''.join(result)
    
    content = replace_enums(content)
    
    # 6. Remove backtick quoting, replace with double quotes for reserved words
    # Reserved words in the schema: `group`
    content = content.replace('`', '"')
    
    # 7. Remove ON UPDATE CURRENT_TIMESTAMP
    content = re.sub(
        r'\s+ON\s+UPDATE\s+CURRENT_TIMESTAMP',
        '',
        content,
        flags=re.IGNORECASE
    )
    
    # 7b. Remove MySQL COMMENT syntax from column definitions
    # Pattern:  type VARCHAR(20) COMMENT 'text' → type VARCHAR(20)
    content = re.sub(
        r"\s+COMMENT\s+'[^']*'",
        '',
        content
    )
    
    # 8. Replace UNIQUE KEY with UNIQUE
    content = re.sub(r'\bUNIQUE\s+KEY\b', 'UNIQUE', content)
    
    # 9. Fix PostgreSQL UNIQUE constraint syntax.
    # MySQL: UNIQUE constraint_name (col1, col2)
    # PG:    CONSTRAINT constraint_name UNIQUE (col1, col2)
    # Pattern: UNIQUE <name> (<columns>)
    # Must be at start of line (table-level constraint)
    content = re.sub(
        r'^\s*UNIQUE\s+(\w+)\s+\(',
        r'CONSTRAINT \1 UNIQUE (',
        content,
        flags=re.MULTILINE
    )
    
    # 10. Replace INDEX (standalone) - no longer needed since step 12 handles inline INDEX
    
    # 11. Handle TIMESTAMP NULL - in PG, TIMESTAMP allows NULL by default
    # Just remove explicit NULL after TIMESTAMP
    content = re.sub(r'\bTIMESTAMP\s+NULL\b', 'TIMESTAMP', content)
    
    # 12. Fix DECIMAL without precision - add default
    # No change needed for our files
    
    # 13. Replace MySQL YEAR type with INTEGER (PostgreSQL doesn't have YEAR)
    content = re.sub(r'\bYEAR\b', 'INTEGER', content)
    
    # 14. Replace MySQL DATETIME with TIMESTAMP (PostgreSQL doesn't have DATETIME)
    content = re.sub(r'\bDATETIME\b', 'TIMESTAMP', content)
    
    # 15. Replace MySQL LONGTEXT/MEDIUMTEXT with TEXT (PostgreSQL uses TEXT for all)
    content = re.sub(r'\bLONGTEXT\b', 'TEXT', content)
    content = re.sub(r'\bMEDIUMTEXT\b', 'TEXT', content)
    
    # 15. For down files: Replace MySQL DROP FOREIGN KEY IF EXISTS with PG DROP CONSTRAINT IF EXISTS
    content = content.replace(
        'DROP FOREIGN KEY IF EXISTS',
        'DROP CONSTRAINT IF EXISTS'
    )
    content = content.replace(
        'DROP FOREIGN KEY',
        'DROP CONSTRAINT'
    )

    # 12. Extract inline INDEX definitions and move them after CREATE TABLE
    # PostgreSQL does NOT support INDEX inside CREATE TABLE
    # Pattern: split each CREATE TABLE block, remove inline INDEX, add CREATE INDEX after
    lines = content.split('\n')
    result_lines = []
    current_table_name = None
    in_create_table = False
    brace_depth = 0
    accumulated_lines = []
    index_statements = []
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        
        # Detect CREATE TABLE start
        create_match = re.match(r'^CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(\w+(\.\w+)?)', stripped, re.IGNORECASE)
        if create_match and not in_create_table:
            in_create_table = True
            current_table_name = create_match.group(1)
            accumulated_lines = [line]
            brace_depth = stripped.count('(') - stripped.count(')')
            index_statements = []
            continue
        
        if in_create_table:
            accumulated_lines.append(line)
            brace_depth += stripped.count('(') - stripped.count(')')
            
            # Check if this line is an inline INDEX
            idx_match = re.match(r'^\s*INDEX\s+', stripped, re.IGNORECASE)
            if idx_match:
                # Extract table name from the INDEX line to handle compound names
                # Pattern: INDEX idx_name (col1, col2) 
                index_name_match = re.search(r'INDEX\s+(\w+)\s+\(([^)]+)\)', stripped, re.IGNORECASE)
                if index_name_match and current_table_name:
                    index_name = index_name_match.group(1)
                    index_columns = index_name_match.group(2)
                    index_sql = f"CREATE INDEX IF NOT EXISTS {index_name} ON {current_table_name} ({index_columns});"
                    index_statements.append(index_sql)
                # Remove this line from accumulated content
                accumulated_lines.pop()
            
            # Check if we've closed the CREATE TABLE
            if brace_depth <= 0 and len(accumulated_lines) > 1:
                # End of CREATE TABLE
                in_create_table = False
                # Write accumulated CREATE TABLE (minus any trailing comma after last column)
                table_content = '\n'.join(accumulated_lines)
                # Remove trailing comma before closing paren
                table_content = re.sub(r',\s*\n\s*\);', '\n);', table_content)
                result_lines.append(table_content)
                # Write CREATE INDEX statements after the table
                for idx_sql in index_statements:
                    result_lines.append('')
                    result_lines.append(idx_sql)
                current_table_name = None
                continue
        else:
            result_lines.append(line)
    
    # If we never entered a CREATE TABLE block, use original lines
    if not any('CREATE TABLE' in l for l in lines):
        result_lines = lines
    
    # If in_create_table is still True (unclosed block), fall back
    if in_create_table:
        result_lines = lines
    
    return '\n'.join(result_lines)

def main():
    # Process .sql files (both up and down)
    for filename in os.listdir(MYSQL_DIR):
        if not filename.endswith('.sql'):
            continue
        if filename.endswith('.down.sql'):
            # Down files: copy and convert MySQL-only syntax to PG
            mysql_path = os.path.join(MYSQL_DIR, filename)
            pg_path = os.path.join(PG_DIR, filename)
            
            with open(mysql_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            converted = convert_sql(content)
            
            with open(pg_path, 'w', encoding='utf-8') as f:
                f.write(converted)
            print(f"Converted (down): {filename}")
            continue
        
        mysql_path = os.path.join(MYSQL_DIR, filename)
        pg_path = os.path.join(PG_DIR, filename)
        
        with open(mysql_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        converted = convert_sql(content)
        
        with open(pg_path, 'w', encoding='utf-8') as f:
            f.write(converted)
        
        print(f"Converted: {filename}")

    # Count files
    mysql_count = len([f for f in os.listdir(MYSQL_DIR) if f.endswith('.sql')])
    pg_count = len([f for f in os.listdir(PG_DIR) if f.endswith('.sql')])
    
    print(f"\nDone! MySQL dir: {mysql_count} files, PostgreSQL dir: {pg_count} files")

if __name__ == '__main__':
    main()
