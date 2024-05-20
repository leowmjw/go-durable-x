# go-durable-x
Compare Durable Execution + Promise Examples 

## Scenarios

### DB Blue Green

### AWS PrivateLink Updater

### AWS DocumentDB Temporary Instance Boost

### Onboard Multi-tenant to DB via ShardingSphere

### Breakglass

### Telegram Note Processor

- Pull latest unproccessed item from Personal Telegram Note Queue (How? SDKs?)
- Normal articles or Twit; gather them into per day Summary; summarized using Kagi Universal Summarizer
- Feed it into dynamic processor depending on type + rules: e.g. PDF split + further process

### PDF File Processor

- Split file by chapter for inpoupt PDF; use LLM to recognize via visual few-shot .. as plan for review
- Plan as YAML; easily adjustable by final human review on a whole
- Feed to Kagi Universal Summarizer; for each chapter/fragment
- Gather it all back together; marking any with failure not able to be addressed for manual intervention
- Have text available nicely as Git shard ..
