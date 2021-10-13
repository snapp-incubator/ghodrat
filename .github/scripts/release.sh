# !bin/sh

project_name=$1
project_commit_sha=$2
commands=$3
registry=$4
username=$5
password=$6

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
