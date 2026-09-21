import sharp from 'sharp';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const inputDir = path.resolve(__dirname, '../../LOKARI LOGO ASSETS');
const outputDir = path.resolve(__dirname, '../static');

async function convertLogos() {
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  // Favicon (just the icon, Asset 10)
  const faviconPath = path.join(inputDir, 'Asset 10@1000x.png');
  if (fs.existsSync(faviconPath)) {
    await sharp(faviconPath)
      .resize(256, 256, { fit: 'contain', background: { r: 0, g: 0, b: 0, alpha: 0 } })
      .webp({ quality: 90 })
      .toFile(path.join(outputDir, 'favicon.webp'));
    console.log('Created favicon.webp');
  }

  // Navbar logo (icon, Asset 10)
  const logoPath = path.join(inputDir, 'Asset 10@1000x.png');
  if (fs.existsSync(logoPath)) {
    await sharp(logoPath)
      .resize(500, null, { withoutEnlargement: true })
      .webp({ quality: 90 })
      .toFile(path.join(outputDir, 'logo.webp'));
    console.log('Created logo.webp');
  }
}

convertLogos().catch(console.error);
