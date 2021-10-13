# !bin/sh

project_name=$1
project_commit_sha=$2
commands="janus"

registry=$3
username=$4
password=$5

# logs docker information
docker info

# login to provided docker registry
docker login -u $username -p $password $registry
    
for command in $commands; do
    # load and unarchive docker image tar file 
    docker load --input $project_name-"$command"-$project_commit_sha.tar

    # push the image to the registry
    docker push $project_name-$command:$project_commit_sha
done
