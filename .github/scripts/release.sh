# !bin/sh

project_name
project_commit_sha

registry
username
password

# logs docker information
docker info

# login to provided docker registry
docker login -u $username -p $password $registry
    
for command in $COMMANDS; do
    # load and unarchive docker image tar file 
    docker load --input $project_name-"$command"-$project_commit_sha.tar

    # push the image to the registry
    docker push $project_name-$command:$project_commit_sha
done
