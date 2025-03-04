# Usar a imagem oficial do Python como base
FROM python:3.11-slim

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar arquivos necessários para o contêiner
COPY requirements.txt requirements.txt
COPY src/ src/

# Instalar dependências
RUN pip install --no-cache-dir -r requirements.txt

# Expor a porta para comunicação via UDP (ajuste conforme necessário)
EXPOSE 20777/udp

# Definir a variável de ambiente para Kafka (caso necessário)
ENV KAFKA_BROKER="kafka:9092"

# Comando para rodar a aplicação
CMD ["python", "src/main.py"]
