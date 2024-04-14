# Utiliza la imagen base oficial de Go
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del directorio actual al directorio de trabajo en el contenedor
COPY . .

# Compila la aplicación Go dentro del contenedor
RUN go build cmd/main.go

# Expone el puerto en el que se ejecutará la aplicación
EXPOSE 8080

# Ejecuta la aplicación cuando se inicie el contenedor
CMD ["./main"]
