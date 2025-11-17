#!/bin/bash
set -e

printf "Warning: the install script needs to be updated, some things may be broken, but it may work\n" 
printf "Continue?\n"
printf "[y/n]\n"
read cont

if [[ "$cont" == "n" ]]; then
  exit 0
else
  printf "continuing"
fi 

projectName="free-games-checker"
projectAuthor="Supraboy981322"

printf "Would you like to fetch binaries (fastest and safest), compile the code locally?\n"
printf "b [binaries]  or  c [compile]\n"
read binOrComp

printf "checking for dependencies\n"
for dep in "jq" "bash" "bzip2" "tar"; do
    if [ -z $(command -v ${dep}) ]; then
         missingDeps+=" - '${dep}'\n"
    fi
done

if [[ ! -z "${missingDeps}" ]]; then
    printf "The ${free-games-checker} install is missing the following dependencies:\n"
    printf "${missingDeps}"
    exit 1
fi

function setSettingsJSON() {
    printf "${projectName} requires 3 Discord webhooks per store.\n"
    printf "If you do not know how to get Discord webhooks, the internet is a great place to ask.\n"
    printf "Continue?\n"
    printf "[y/n]\n"
    read cont
    ls
    if [[ "$cont" == "y" ]]; then
        for store in "Steam" "GOG"; do
            case "${store}" in
                "Steam")
                    storeNum=0
                    ;;
                "GOG")
                    storeNum=1
                    ;;
                "Itch")
                    storeNum=2
                    ;;
                "Epic")
                    storeNum=3
                    ;;
                *)
                    printf "ERR! INVALID STORE ENTERED!"
                    printf "THIS LIKELY MEANS SOMEONE'S DOING SOMETHING NASTY!"
                    exit 1
                    ;;
            esac
            printf "Now, you will need webhooks to use for ${store}\n"
            for num in $(seq 0 2); do
                printf "Please enter webhook $((num+1)) for ${store}\n"
                read webhook
                jq ".[${storeNum}].webhooks.[${num}] = \"${webhook}\"" settings.json > tmp-settings.json
                rm settings.json
                mv tmp-settings.json settings.json
            done
        done
        printf "settings complete"
    else
        printf "exiting...\n"
        exit 0
    fi
}

case "$binOrComp" in
    "b")
        printf "creating 'free-games-checker' directory\n"
        mkdir free-games-checker
        cd free-games-checker
        printf "fetching '.bzip2' tarball\n"
        wget "https://raw.githubusercontent.com/${projectAuthor}/${projectName}/main/build/package/free-games-checker.tar.bz2" &> /dev/null
        printf "extracting '.bzip2 tarball\n"
        tar -xjf free-games-checker.tar.bz2 &> /dev/null
        printf "cleanup\n"
        rm free-games-checker.tar.bz2
        setSettingsJSON
        ;;
    *)
        printf "unrecognized option, exiting\n"
        ;;
esac

