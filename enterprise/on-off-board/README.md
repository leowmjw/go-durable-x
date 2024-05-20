# On / Off Boarding

## Scenarios

- All Enterprises have employees coming on board; and might be around for very short or very long time
- It is always very tedious; and every company is different .. let's make it all seamless to add on plugins + workflows 
- HOWTO handle all lifecycles
	- Ad-hoc actions like add new tools
	- Human in the loop actions
	- Newer version of workflow with the new tools
	- Able to automate as much backported for existing users; into new tools
	- LLM agents in the loop ..?
	- Off-boarding automatically; following the depency tree; e.g. remove Okta last
	- Auto-reverse permissions; before making it permanent but locking username
	- Merging new company + migratiing over ..

## Use Case

- Added manually into AWS IAM via Console; mark human tasks completed
- Added automatically into PagerDuty via API
- Change workflow to handle 
- Generates Terraform; to be applied manually; human signals
- Change workflow so that it now gets based on Github workflow from Atlantis?
- Automatic adding into Okta; has an audit trail
- Migarte over from AWS IAM into AWS Identity Center

## Unit tests

- Do TDD for above scenarios; to be added incrementally ...
