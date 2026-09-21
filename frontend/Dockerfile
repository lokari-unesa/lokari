FROM node:20-alpine

# Set working directory
WORKDIR /app

# Salin file konfigurasi package
COPY package.json package-lock.json* ./

# Install dependensi
RUN npm install

# Salin seluruh kode frontend
COPY . .

# Ekspos port yang digunakan oleh SvelteKit (sesuai VITE_PORT)
EXPOSE 5180

# Jalankan server dalam mode development dan pastikan terekspos ke luar container (0.0.0.0)
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
