# AWS / DocuDB Boost

## Scenario

- DocuDB cluster has the smallest tony instance
- Scheduled time during US Workday will scale instance to small instance
- Receiving signal after CPU <5% for more than 5 minutes
- DocuDB cluster is shrunk back down to original tony instance
