# !/bin/bash

APP_NAME=kafka-producer
AWS_REGION="eu-west-1"
REGISTRY_ID=$(aws ecr describe-registry --output text --query 'registryId' --region $AWS_REGION)
REGISTRY="${REGISTRY_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
REPO_PREFIX=infra
TAG=$(git rev-parse --short HEAD --)

# Stopping the script if any command fails
set -e

# === Build and push Docker image === 
aws ecr describe-repositories --region $AWS_REGION --repository-names ${REPO_PREFIX}/${APP_NAME} || aws ecr create-repository --repository-name ${REPO_PREFIX}/${APP_NAME} --region $AWS_REGION 
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $REGISTRY
echo "==== Building GoLang binary ===="
go mod download

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm64
go build -buildvcs=false -o ${APP_NAME}-arm64 .

export GOARCH=amd64
go build -buildvcs=false -o ${APP_NAME}-amd64 .
# docker buildx build --push --platform linux/amd64,linux/arm64 -t ${REGISTRY}/${REPO_PREFIX}/${APP_NAME}:${TAG} .
docker buildx build --push --platform linux/arm64 -t ${REGISTRY}/${REPO_PREFIX}/${APP_NAME}:${TAG} -t ${REGISTRY}/${REPO_PREFIX}/${APP_NAME}:latest .

rm -f ${APP_NAME}-amd64 ${APP_NAME}-arm64
