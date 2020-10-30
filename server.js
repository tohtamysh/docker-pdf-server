const puppeteer = require('puppeteer');

const express = require('express');

const bodyParser = require('body-parser');

const urlencodedParser = bodyParser.urlencoded({ limit: '50mb', extended: false });

const server = express();

function getUrl(url, props, res) {
  (async () => {
    const browser = await puppeteer.launch({
      headless: true,
      args: ['--disable-dev-shm-usage', '--no-sandbox', '--disable-setuid-sandbox', '--allow-file-access-from-files', '--enable-local-file-accesses', '--disable-web-security']
    });

    const page = await browser.newPage();

    await page.goto(url, { waitUntil: 'networkidle0' })

    const pdfBuffer = await page.pdf({
      format: 'A4',
      landscape: props.landscape
    });

    res.writeHead(200, {
      'Content-Type': 'application/pdf',
      'Content-Length': pdfBuffer.length
    });
    res.end(pdfBuffer);

    await browser.close();
  })()
}

function getHtml(html, props, res) {
  (async () => {
    const browser = await puppeteer.launch({
      headless: true,
      args: ['--disable-dev-shm-usage', '--no-sandbox', '--disable-setuid-sandbox', '--allow-file-access-from-files', '--enable-local-file-accesses', '--disable-web-security']
    });

    const page = await browser.newPage();

    await page.setContent(html, { waitUntil: 'networkidle0' });

    const pdfBuffer = await page.pdf({
      format: 'A4',
      landscape: props.landscape
    });

    res.writeHead(200, {
      'Content-Type': 'application/pdf',
      'Content-Length': pdfBuffer.length
    });
    res.end(pdfBuffer);

    await browser.close();
  })()
}

server.post('/', urlencodedParser, function (req, res) {
  const url = req.body.url ? req.body.url : null;

  const html = req.body.html ? req.body.html : null;

  const landscape = new Boolean(req.body.landscape);

  const props = {
    landscape: landscape,
  }

  if (!url && !html) return res.sendStatus(400);

  if (url) {
    getUrl(url, props, res);
  }

  if (html) {
    getHtml(html, props, res);
  }
});

server.listen(8000, () => {
  console.log('Server started!');
});