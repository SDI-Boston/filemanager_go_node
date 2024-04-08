FROM ubuntu:latest

# Actualizar repositorios e instalar dependencias
RUN apt-get update \
    && apt-get install -y \
        golang nfs-common\
    && rm -rf /var/lib/apt/lists/*

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el código fuente de la aplicación
COPY . .

EXPOSE 5000

# Comando para ejecutar la aplicación
CMD ["go", "run", "main.go"]
