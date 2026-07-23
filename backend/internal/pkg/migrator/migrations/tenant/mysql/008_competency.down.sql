-- Down Migration: 008_competency
-- Drop FK constraint first (added by 008) before dropping referenced table
ALTER TABLE job_family_competencies DROP FOREIGN KEY fk_jfc_competency;

DROP TABLE IF EXISTS competency_score_details;
DROP TABLE IF EXISTS competency_scores;
DROP TABLE IF EXISTS competency_event_targets;
DROP TABLE IF EXISTS competency_events;
DROP TABLE IF EXISTS competency_values;
DROP TABLE IF EXISTS competence_values;
DROP TABLE IF EXISTS competencies;
