#!/bin/bash

# ========================================
# Build and Push zpwoot to Docker Hub
# ========================================

set -e

# Cores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Variáveis
DOCKER_USERNAME="${DOCKER_USERNAME:-}"
IMAGE_NAME="zpwoot"
VERSION="${VERSION:-latest}"

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Build and Push zpwoot to Docker Hub${NC}"
echo -e "${YELLOW}========================================${NC}"

# Verificar se o usuário do Docker Hub foi fornecido
if [ -z "$DOCKER_USERNAME" ]; then
    echo -e "${YELLOW}Digite seu username do Docker Hub:${NC}"
    read -r DOCKER_USERNAME
fi

if [ -z "$DOCKER_USERNAME" ]; then
    echo -e "${RED}❌ Username do Docker Hub é obrigatório!${NC}"
    exit 1
fi

FULL_IMAGE_NAME="${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

echo -e "${GREEN}📦 Imagem: ${FULL_IMAGE_NAME}${NC}"
echo ""

# Login no Docker Hub
echo -e "${YELLOW}🔐 Fazendo login no Docker Hub...${NC}"
docker login

# Build da imagem
echo -e "${YELLOW}🔨 Building imagem...${NC}"
docker build -t "${FULL_IMAGE_NAME}" -f Dockerfile .

# Tag latest
if [ "$VERSION" != "latest" ]; then
    echo -e "${YELLOW}🏷️  Criando tag latest...${NC}"
    docker tag "${FULL_IMAGE_NAME}" "${DOCKER_USERNAME}/${IMAGE_NAME}:latest"
fi

# Push para Docker Hub
echo -e "${YELLOW}📤 Pushing para Docker Hub...${NC}"
docker push "${FULL_IMAGE_NAME}"

if [ "$VERSION" != "latest" ]; then
    docker push "${DOCKER_USERNAME}/${IMAGE_NAME}:latest"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Imagem publicada com sucesso!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Imagem: ${FULL_IMAGE_NAME}${NC}"
echo ""
echo -e "${YELLOW}Para usar no stack.yml, atualize:${NC}"
echo -e "  image: ${FULL_IMAGE_NAME}"
echo ""

