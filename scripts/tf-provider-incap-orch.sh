#!/bin/bash

# AUTHOR: Pablo S. Martinez 
# STAKEHOLDERS: Imperva TF Provider developers
# DESCRIPTION: Following script is an automation on top of provider make commands 
# and it aims to automate and wrap up manual iterative configurations
# Run ./tf-provider-incap-orch.sh to see the menu options and script usage, in brief
# script currently supports the following capabilities:

#Install
#1. Clones provider git repo or pull latest
#2. Creates plugin directory
#3. Builds provider and copies binary to plugin destination folder
#4. Creates go mapper to plugin destination
#5. Creates terraform directory and intial configuration files 

#Update
#1. Builds provider and copies binary to plugin destination folder
#2. Updates provider version in terrfaorm main.tf file

#Clean
#1. Clean terraform directory
#2. Clean plugin directory
#3. Clean provider binary and downloaded go dependencies

#Unit Tests
#1. Execute unit tests

#Acceptance Tets
#1. Execute acceptance tests

USER=`whoami`
ROOT=/Users/$USER
GOPATH=$ROOT/workspace/go
PROVIDER_PLUGIN=$GOPATH/plugins/registry.terraform.io/terraform-providers/incapsula
PROVIDER_GIT=https://github.com/imperva/terraform-provider-incapsula.git
PROVIDER_GIT_LOCAL=$GOPATH/src/github.com/terraform-providers/terraform-provider-incapsula
TERRAFORM=$ROOT/workspace/terraform
MAIN_TF=$TERRAFORM/main.tf
VARS_TF=$TERRAFORM/terraform.tfvars
I_FLAG="false"
B_FLAG="false"
T_FLAG="false"
A_FLAG="false"
C_FLAG="false"

export INCAPSULA_API_ID=$2
export INCAPSULA_API_KEY=$3
export GO111MODULE="on" 

print_usage() {
  printf "\n****************************\n
Usage: ./tf-provider-incap-orch.sh <option-character> <optional-API-ID> <optional-API-Key> \n
Options:\n
 -i Install, usage example: ./tf-provider-incap-orch.sh -i \"12345\" \"a1fsa2f-fas24fsaf\" \n
 -b Build, usage example: ./tf-provider-incap-orch.sh -b \n 
 -t Unit Tests, usage example: /tf-provider-incap-orch.sh -t \"12345\" \"a1fsa2f-fas24fsaf\"  \n 
 -a Acceptance Tests, usage example: /tf-provider-incap-orch.sh -a \"12345\" \"a1fsa2f-fas24fsaf\"  \n 
 -c Clean, usage example: ./tf-provider-incap-orch.sh -c \n\n****************************\n"
}

validation(){
    log_entry ${FUNCNAME[0]}
    if [ ! -f "`which brew`" ]; then log_error ${FUNCNAME[0]} 'Brew is not installed'; exit 1; fi;
}

install(){      
    log_entry ${FUNCNAME[0]}
    validation

    if [ ! -f "`which go`" ]; 
    then 
        log_info "${FUNCNAME[0]}" "Installing Golang..."
        brew install go; 
    else
        log_info "${FUNCNAME[0]}" "Golang detected"
    fi;

    if [ ! -f "`which git`" ];   
    then 
        log_info "${FUNCNAME[0]}" "Installing GIT..."
        brew install git; 
    else
        log_info "${FUNCNAME[0]}" "GIT detected"
    fi;

    if [ ! -f "`which terraform`" ];   
    then 
        log_info "${FUNCNAME[0]}" "Installing Terraform local client..."
        brew tap hashicorp/tap
        brew install hashicorp/tap/terraform
    else
        log_info "${FUNCNAME[0]}" "Terraform detected"
    fi;

    log_info "${FUNCNAME[0]}" "Login user is ${USER}"

    mkdir -p $GOPATH
    chmod -R 777 $GOPATH
    mkdir -p $TERRAFORM
    chmod -R 777 $TERRAFORM

    log_info "${FUNCNAME[0]}" "Terraform GIT path $PROVIDER_GIT_LOCAL"

    if [ -d "$PROVIDER_GIT_LOCAL" ]; 
    then 
        log_info "${FUNCNAME[0]}" "Pulling repo terraform-provider-incapsula"
        git --git-dir=$PROVIDER_GIT_LOCAL/.git config core.fileMode false
        git --git-dir=$PROVIDER_GIT_LOCAL/.git pull
    else 
        log_info "${FUNCNAME[0]}" "Cloning repo terraform-provider-incapsula"
        git clone $PROVIDER_GIT $PROVIDER_GIT_LOCAL
    fi

    log_info "${FUNCNAME[0]}" "Creating ~/.terraformrc file"
    
    echo "provider_installation {
  filesystem_mirror {
    path    = \"/Users/$USER/workspace/go/plugins\"
    include = [\"registry.terraform.io/terraform-providers/incapsula\"]
  }
}" > ~/.terraformrc

    PROVIDER_VERSION=`git --git-dir=$PROVIDER_GIT_LOCAL/.git tag | tail -1 | cut -c2-`
    log_info "${FUNCNAME[0]}" "Git provider version $PROVIDER_VERSION"

    mkdir -p $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/
    log_info "${FUNCNAME[0]}" "Provider binary path $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/"

    log_info "${FUNCNAME[0]}" "Building local provider..."
    make -C $PROVIDER_GIT_LOCAL fmt && make -C $PROVIDER_GIT_LOCAL build    
    rm -rf $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/ ||:
    cp -r $PROVIDER_GIT_LOCAL/ $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/

    if [ ! -f "$MAIN_TF" ]; 
    then 
        log_info "${FUNCNAME[0]}" "Created file $MAIN_TF"
        echo "terraform {
  required_providers {
    incapsula = {
      source = \"terraform-providers/incapsula\"
      version = \"$PROVIDER_VERSION\"
    }
  }
}
 
variable \"incapsula_api_id\" {
  type        = number
  description = \"API ID\"
}
 
variable \"incapsula_api_key\" {
  type        = string
  description = \"API KEY\"
}

provider \"incapsula\" {
  api_id = var.incapsula_api_id
  api_key = var.incapsula_api_key
}" > $MAIN_TF

        chmod 777 $MAIN_TF
    fi;

    if [ ! -f "$VARS_TF" ]; 
    then 
        log_info "${FUNCNAME[0]}" "Created file $VARS_TF"
        echo "incapsula_api_key = \"$INCAPSULA_API_KEY\"
incapsula_api_id = $INCAPSULA_API_ID" > $VARS_TF 
        chmod 777 $VARS_TF
    fi;
}

build(){
    log_entry ${FUNCNAME[0]}
    make -C $PROVIDER_GIT_LOCAL fmt && make -C $PROVIDER_GIT_LOCAL build

    PROVIDER_VERSION=`grep "VERSION=" $PROVIDER_GIT_LOCAL/GNUmakefile | cut -d "=" -f2`
    log_info "${FUNCNAME[0]}" "Local provider version $PROVIDER_VERSION"
    rm -rf $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/ ||:
    cp -r $PROVIDER_GIT_LOCAL/ $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/ 
    log_info "${FUNCNAME[0]}" "Provider binary path $PROVIDER_PLUGIN/$PROVIDER_VERSION/darwin_amd64/"

    sed -i '' "s/version = .*/version = \"${PROVIDER_VERSION}\"/" $MAIN_TF
}

unit(){
    log_entry ${FUNCNAME[0]}
    cd $PROVIDER_GIT_LOCAL/incapsula && go test
}

acceptance(){
    log_entry ${FUNCNAME[0]}
    cd $PROVIDER_GIT_LOCAL && make testacc
}

clean(){
    log_entry ${FUNCNAME[0]}
    rm -rf $PROVIDER_PLUGIN/ ||:
    rm -rf $TERRAFORM/ ||:
    make -C $PROVIDER_GIT_LOCAL clean
}

log(){
    if [ -z "$3" ]
    then
        echo `date` - $1 - method: $2\(\)
    else
        echo `date` - $1 - method: $2\(\) message: $3
    fi;
}

log_info(){
    log "$MSG_INFO" "$1" "$2" 
}

log_error(){
    log "$MSG_ERROR" "$1" "$2" 
}

log_entry(){
    log "$MSG_ENTRY" "$1"
}

# main

while getopts 'ibtac' flag; do
  case "${flag}" in
    i) I_FLAG='true' ;;
    b) B_FLAG='true' ;;
    t) T_FLAG='true' ;;
    a) A_FLAG='true' ;;
    c) C_FLAG='true' ;;
    *) print_usage
       exit 1 ;;
  esac
done

if [ "$I_FLAG" == 'true' ]
then 
    install; 
elif [ "$B_FLAG" == 'true' ]
then 
    build;
elif [ "$T_FLAG" == 'true' ]
then 
    unit;
elif [ "$A_FLAG" == 'true' ]
then 
    acceptance;
elif [ "$C_FLAG" == 'true' ]
then 
    clean;
else 
    print_usage;
fi;
