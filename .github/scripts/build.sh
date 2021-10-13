# !bin/sh

project_name=$1
project_commit_sha=$2
commands="janus"

export CURRENT_DATETIME=$(TZ=Asia/Tehran date '+%FT%T')

render() {
    sedStr="s!%%COMMAND%%!$command!g;"
    sed -r "$sedStr" $1
}

buildDirectory="./build"

for command in $commands; do
    # generate dockerfile for following builds
    render $buildDirectory/template.Dockerfile > $buildDirectory/$command.Dockerfile

    # build image based on Dockerfile and build-arguments
    docker build \
        --build-arg BUILD_DATE=$CURRENT_DATETIME \
        -t $project_name-$command:$project_commit_sha \
        -f $buildDirectory/$command.Dockerfile .

    # archive builded image as a tar file
    docker save \
        -o $project_name-"$command"-$project_commit_sha.tar \
        $project_name-$command:$project_commit_sha
done
