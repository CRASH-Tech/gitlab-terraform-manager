**The tool for backup and restore gitlab terraform states**

**Build how to:**
```
git clone https://github.com/CRASH-Tech/gitlab-terraform-manager.git
cd gitlab-terraform-manager/cmd/gitlab-terraform-manager/
go build
```
**Just set config vars for access to gilab and relax**

```
export GITLAB_URL="https://gitlab.com"
export GITLAB_TOKEN="8vJxe8UDpN5mHYG"
export GITLAB_PROJECT_ID="100"
export GITLAB_PROJECT_PATH="infra/terraform"
```

**Get list of state files**

```./gitlab-terraform-manager list```

**Remove state file**

```./gitlab-terraform-manager remove example_tf_state```

**Save state file, filename will be generated automatically with timestamp. DO NOT rename it!**

```./gitlab-terraform-manager save example_tf_state```

**Save sll state files to folder**

```./gitlab-terraform-manager saveall data/```

**Restore all state files from folder. If we have more than one file with same name, will be restored only state with latest timestamp**

```./gitlab-terraform-manager restoreall data/```
