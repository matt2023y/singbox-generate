// server.js
const http = require('http');
const { URL } = require('url');

const PORT = 3000;

const server = http.createServer(async (req, res) => {
    try {
        // 收集请求体
        const chunks = [];
        for await (const chunk of req) chunks.push(chunk);
        const rawBody = Buffer.concat(chunks).toString();

        // 解析 URL 和查询参数
        const fullUrl = `http://${req.headers.host || 'localhost'}${req.url}`;
        const urlObj = new URL(fullUrl);
        const query = {};
        for (const [k, v] of urlObj.searchParams) query[k] = v;

        // 构造请求信息对象
        const info = {
            timestamp: new Date().toISOString(),
            method: req.method,
            url: urlObj.pathname,
            fullUrl: fullUrl,
            query: query,
            headers: req.headers,
            rawBody: rawBody,
        };

        // 打印到控制台（格式化）
        console.log('--- Incoming Request ---');
        console.log(JSON.stringify(info, null, 2));
        console.log('------------------------');

        // 返回 JSON 给客户端
        const body = JSON.stringify(info, null, 2);
        res.writeHead(200, {
            'Content-Type': 'application/json; charset=utf-8',
            'Content-Length': Buffer.byteLength(body),
        });
        res.end(body);
    } catch (err) {
        console.error('Error handling request:', err);
        res.writeHead(500, { 'Content-Type': 'text/plain; charset=utf-8' });
        res.end('Internal Server Error');
    }
});

server.listen(PORT, () => {
    console.log(`Server listening on http://localhost:${PORT}`);
});
