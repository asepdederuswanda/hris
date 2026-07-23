-- Down Migration: 010_approval
DROP TABLE IF EXISTS approval_tasks;
DROP TABLE IF EXISTS approval_actions;
DROP TABLE IF EXISTS approval_instances;
DROP TABLE IF EXISTS approval_flow_steps;
DROP TABLE IF EXISTS approval_flows;
