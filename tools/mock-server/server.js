#!/usr/bin/env node
/**
 * CastBrick Canonical Mock API Server
 * Lightweight, zero-dependency HTTP server for automated testing of CastBrick SDKs.
 */

import http from 'node:http';
import { URL, fileURLToPath } from 'node:url';
import fs from 'node:fs';
import path from 'node:path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Load canonical OpenAPI specification
let openApiSpec = null;
const possibleSpecPaths = [
  path.resolve(__dirname, '../specs/openapi.json'),
  path.resolve(__dirname, '../../tools/specs/openapi.json'),
  path.resolve(__dirname, './specs/openapi.json'),
  path.resolve(process.cwd(), 'tools/specs/openapi.json'),
];

for (const p of possibleSpecPaths) {
  if (fs.existsSync(p)) {
    try {
      openApiSpec = JSON.parse(fs.readFileSync(p, 'utf-8'));
      break;
    } catch {}
  }
}

if (openApiSpec) {
  console.log(`[CASTBRICK-MOCK] OpenAPI contract loaded: "${openApiSpec.info?.title} v${openApiSpec.info?.version}"`);
} else {
  console.warn('[CASTBRICK-MOCK] Warning: OpenAPI spec not found. Running without strict schema validation.');
}

function validateOpenApiPayload(schemaName, body) {
  if (!openApiSpec || !openApiSpec.components || !openApiSpec.components.schemas) {
    return null;
  }
  const schema = openApiSpec.components.schemas[schemaName];
  if (!schema) return null;

  if (schema.required) {
    for (const field of schema.required) {
      if (body[field] === undefined || body[field] === null || body[field] === '') {
        return `OpenAPI Contract Error: Missing required field '${field}' (schema: ${schemaName})`;
      }
    }
  }

  if (schema.properties) {
    for (const [prop, propSchema] of Object.entries(schema.properties)) {
      if (body[prop] !== undefined && body[prop] !== null) {
        if (propSchema.type === 'array' && !Array.isArray(body[prop])) {
          return `OpenAPI Contract Error: Field '${prop}' must be an array (schema: ${schemaName})`;
        }
        if (propSchema.type === 'string' && typeof body[prop] !== 'string') {
          return `OpenAPI Contract Error: Field '${prop}' must be a string (schema: ${schemaName})`;
        }
        if (propSchema.type === 'boolean' && typeof body[prop] !== 'boolean') {
          return `OpenAPI Contract Error: Field '${prop}' must be a boolean (schema: ${schemaName})`;
        }
        if (propSchema.type === 'number' && typeof body[prop] !== 'number') {
          return `OpenAPI Contract Error: Field '${prop}' must be a number (schema: ${schemaName})`;
        }
      }
    }
  }
  return null;
}

const PORT = parseInt(process.env.PORT || process.argv[2] || '8787', 10);
const HOST = process.env.HOST || '127.0.0.1';

function sendJson(res, statusCode, data) {
  const body = JSON.stringify(data);
  res.writeHead(statusCode, {
    'Content-Type': 'application/json',
    'Content-Length': Buffer.byteLength(body),
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Headers': 'Authorization, Content-Type, Accept',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
  });
  res.end(body);
}

function sendNoContent(res) {
  res.writeHead(204, {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Headers': 'Authorization, Content-Type, Accept',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
  });
  res.end();
}

const server = http.createServer(async (req, res) => {
  // CORS Preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204, {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Headers': 'Authorization, Content-Type, Accept',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    });
    return res.end();
  }

  const parsedUrl = new URL(req.url, `http://${req.headers.host || 'localhost'}`);
  let pathname = parsedUrl.pathname;

  // Normalize path (strip /v1 prefix if present)
  if (pathname.startsWith('/v1/')) {
    pathname = pathname.substring(3);
  } else if (pathname === '/v1') {
    pathname = '/';
  }

  // Health check endpoint (public)
  if (pathname === '/health' || pathname === '/') {
    return sendJson(res, 200, {
      status: 'ok',
      service: 'castbrick-mock-api',
      timestamp: new Date().toISOString(),
    });
  }

  // Check Authorization header for all other endpoints
  const authHeader = req.headers['authorization'] || '';
  if (
    !authHeader.startsWith('Bearer ') ||
    authHeader === 'Bearer bad' ||
    authHeader === 'Bearer invalid' ||
    authHeader === 'Bearer invalid_key' ||
    authHeader.length <= 10
  ) {
    return sendJson(res, 401, {
      error: 'Unauthorized',
      message: 'A valid Bearer API key is required.',
    });
  }

  // Read JSON body if POST / PUT
  let body = {};
  if (req.method === 'POST' || req.method === 'PUT') {
    try {
      const chunks = [];
      for await (const chunk of req) {
        chunks.push(chunk);
      }
      const rawBody = Buffer.concat(chunks).toString('utf-8');
      if (rawBody.trim().length > 0) {
        body = JSON.parse(rawBody);
      }
    } catch {
      return sendJson(res, 400, {
        error: 'BadRequest',
        message: 'Invalid JSON payload.',
      });
    }
  }

  // Simulated error trigger
  if (pathname.startsWith('/test/error/')) {
    const code = parseInt(pathname.split('/').pop() || '500', 10);
    return sendJson(res, code, {
      error: 'SimulatedError',
      status: code,
      message: `Triggered simulated error ${code}`,
    });
  }

  // -------------------------------------------------------------
  // SMS Endpoints
  // -------------------------------------------------------------
  if (pathname === '/sms/send' && req.method === 'POST') {
    const contractError = validateOpenApiPayload('SendSmsRequest', body);
    if (contractError) {
      return sendJson(res, 422, {
        error: 'OpenApiContractViolation',
        message: contractError,
        schema: 'SendSmsRequest',
      });
    }

    const recipients = body.recipients || (body.to ? (Array.isArray(body.to) ? body.to : [body.to]) : []);
    const count = recipients.length || 1;
    const now = new Date().toISOString();
    return sendJson(res, 200, {
      id: 'a0000000-0000-0000-0000-000000000001',
      messageId: 'a0000000-0000-0000-0000-000000000001',
      status: 'Queued',
      recipientsCount: count,
      recipientCount: count,
      costAoa: 10,
      timestamp: now,
      createdAt: now,
    });
  }

  if (pathname === '/sms/estimate' && req.method === 'POST') {
    return sendJson(res, 200, {
      partsCount: 1,
      recipientsCount: Array.isArray(body.recipients) ? body.recipients.length : 1,
      estimatedCostAoa: 10,
    });
  }

  if (pathname === '/sms' && req.method === 'GET') {
    const now = new Date().toISOString();
    return sendJson(res, 200, {
      items: [
        {
          id: 'a0000000-0000-0000-0000-000000000001',
          messageId: 'a0000000-0000-0000-0000-000000000001',
          content: 'Hello from CastBrick Mock!',
          message: 'Hello from CastBrick Mock!',
          phone: '+244923000000',
          recipientPhone: '+244923000000',
          status: 'Delivered',
          retryCount: 0,
          createdAt: now,
        },
      ],
      totalCount: 1,
      pageNumber: 1,
      pageSize: 20,
    });
  }

  if (pathname.startsWith('/sms/') && req.method === 'DELETE') {
    return sendNoContent(res);
  }

  if (pathname.startsWith('/sms/') && req.method === 'GET') {
    const id = pathname.replace('/sms/', '');
    return sendJson(res, 200, {
      id: id || 'a0000000-0000-0000-0000-000000000001',
      content: 'Hello from CastBrick Mock!',
      phone: '+244923000000',
      status: 'Delivered',
      createdAt: new Date().toISOString(),
    });
  }

  // -------------------------------------------------------------
  // Contacts & Audience Endpoints
  // -------------------------------------------------------------
  if ((pathname === '/contacts' || pathname === '/audience/contacts') && req.method === 'POST') {
    return sendJson(res, 201, {
      id: 'b0000000-0000-0000-0000-000000000001',
      name: body.name || 'Test Contact',
      phone: body.phone || body.endpoint || '+244923000000',
      phoneNumber: body.phone || body.endpoint || '+244923000000',
      endpoint: body.phone || body.endpoint || '+244923000000',
      tenantId: '00000000-0000-0000-0000-000000000001',
      createdAt: new Date().toISOString(),
    });
  }

  if ((pathname === '/contacts' || pathname === '/audience/contacts') && req.method === 'GET') {
    return sendJson(res, 200, {
      items: [
        {
          id: 'b0000000-0000-0000-0000-000000000001',
          name: 'Jane Doe',
          phone: '+244923000000',
          phoneNumber: '+244923000000',
          endpoint: '+244923000000',
          tenantId: '00000000-0000-0000-0000-000000000001',
          createdAt: new Date().toISOString(),
        },
      ],
      totalCount: 1,
      pageNumber: 1,
      pageSize: 20,
    });
  }

  if ((pathname.startsWith('/contacts/') || pathname.startsWith('/audience/contacts/')) && req.method === 'GET') {
    const id = pathname.split('/').pop();
    return sendJson(res, 200, {
      id: id || 'b0000000-0000-0000-0000-000000000001',
      name: 'Jane Doe',
      phone: '+244923000000',
      phoneNumber: '+244923000000',
      endpoint: '+244923000000',
      tenantId: '00000000-0000-0000-0000-000000000001',
      createdAt: new Date().toISOString(),
    });
  }

  if ((pathname.startsWith('/contacts/') || pathname.startsWith('/audience/contacts/')) && req.method === 'DELETE') {
    return sendNoContent(res);
  }

  if (pathname === '/audience/lists' && req.method === 'GET') {
    return sendJson(res, 200, {
      items: [
        {
          id: 'd0000000-0000-0000-0000-000000000001',
          name: 'Default List',
          contactCount: 10,
          tenantId: '00000000-0000-0000-0000-000000000001',
          createdAt: new Date().toISOString(),
        },
      ],
      totalCount: 1,
      pageNumber: 1,
      pageSize: 20,
    });
  }

  // -------------------------------------------------------------
  // Broadcasts Endpoints
  // -------------------------------------------------------------
  if (pathname === '/broadcasts' && req.method === 'POST') {
    const contractError = validateOpenApiPayload('CreateBroadcastRequest', body);
    if (contractError) {
      return sendJson(res, 422, {
        error: 'OpenApiContractViolation',
        message: contractError,
        schema: 'CreateBroadcastRequest',
      });
    }

    const msg = body.message || body.content || 'Mock broadcast message';
    return sendJson(res, 201, {
      id: 'c0000000-0000-0000-0000-000000000001',
      name: body.name || 'Mock Broadcast Campaign',
      status: 'Scheduled',
      message: msg,
      content: msg,
      totalRecipients: 100,
      tenantId: '00000000-0000-0000-0000-000000000001',
      createdAt: new Date().toISOString(),
    });
  }

  if (pathname === '/broadcasts' && req.method === 'GET') {
    return sendJson(res, 200, {
      items: [
        {
          id: 'c0000000-0000-0000-0000-000000000001',
          name: 'Marketing Launch',
          status: 'Completed',
          message: 'Marketing Launch SMS',
          content: 'Marketing Launch SMS',
          totalRecipients: 50,
          tenantId: '00000000-0000-0000-0000-000000000001',
          createdAt: new Date().toISOString(),
        },
      ],
      totalCount: 1,
      pageNumber: 1,
      pageSize: 20,
    });
  }

  if (pathname.startsWith('/broadcasts/') && req.method === 'GET') {
    const id = pathname.replace('/broadcasts/', '');
    return sendJson(res, 200, {
      id: id || 'c0000000-0000-0000-0000-000000000001',
      name: 'Marketing Launch',
      status: 'Completed',
      message: 'Marketing Launch SMS',
      content: 'Marketing Launch SMS',
      totalRecipients: 50,
      tenantId: '00000000-0000-0000-0000-000000000001',
      createdAt: new Date().toISOString(),
    });
  }

  // -------------------------------------------------------------
  // Push Notifications Endpoints
  // -------------------------------------------------------------
  if (pathname === '/push/tokens' && req.method === 'POST') {
    return sendJson(res, 200, {
      token: 'push_mock_token_abc123',
      expiresAt: new Date(Date.now() + 86400000).toISOString(),
    });
  }

  if (pathname === '/push/publish' && req.method === 'POST') {
    return sendJson(res, 200, {
      messageId: 'push_msg_mock_789',
      channel: body.channel || 'default',
      event: body.event || 'test',
      delivered: true,
      timestamp: new Date().toISOString(),
    });
  }

  // -------------------------------------------------------------
  // Billing / Balance Endpoints
  // -------------------------------------------------------------
  if (pathname === '/billing/balance' && req.method === 'GET') {
    return sendJson(res, 200, {
      balance: 50000.0,
      credits: 5000,
      currency: 'AOA',
      prepaid: true,
      lowBalanceThreshold: 500.0,
      alertEmail: 'admin@castbrick.co',
      lastAlertSentAt: null,
    });
  }

  // -------------------------------------------------------------
  // Fallback 404
  // -------------------------------------------------------------
  return sendJson(res, 404, {
    error: 'NotFound',
    message: `Endpoint ${req.method} ${pathname} not found on mock server.`,
  });
});

server.listen(PORT, HOST, () => {
  console.log(`[CASTBRICK-MOCK] Server listening on http://${HOST}:${PORT}`);
});

function shutdown() {
  console.log('[CASTBRICK-MOCK] Shutting down...');
  server.close(() => {
    process.exit(0);
  });
}

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);
