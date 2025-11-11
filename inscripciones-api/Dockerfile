# Imagen base liviana
FROM alpine:latest

# Crear una carpeta dentro del contenedor
WORKDIR /app

# Copiar el contenido local (incluye inscripciones.json y db/)
COPY . .

# Comando por defecto: listar el contenido (solo para debug)
CMD ["ls", "-R", "/app"]

