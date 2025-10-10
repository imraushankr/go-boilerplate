// import express from "express";

// const app = express();
// const port = 7809;

// app.get("/", (req, res) => {
//   res.send("Hello, this is the docs server!");
// });

// app.listen(port, () => {
//   console.log(`Docs server running at http://localhost:${port}`);
// });

import express from 'express';
import { apiReference } from '@scalar/express-api-reference';

const app = express();
const PORT = process.env.PORT || 3000;

// Custom OpenAPI specification
const openApiSpec = {
  openapi: '3.1.0',
  info: {
    title: 'My Awesome API',
    version: '1.0.0',
    description: 'A sample API with TypeScript and Scalar',
    contact: {
      name: 'API Support',
      email: 'support@example.com',
    },
  },
  servers: [
    {
      url: `http://localhost:${PORT}`,
      description: 'Development server',
    },
  ],
  paths: {
    '/api/v1/users': {
      get: {
        summary: 'Get all users',
        description: 'Retrieve a list of all users',
        tags: ['Users'],
        responses: {
          '200': {
            description: 'Successful response',
            content: {
              'application/json': {
                schema: {
                  type: 'array',
                  items: {
                    type: 'object',
                    properties: {
                      id: { 
                        type: 'integer',
                        example: 1
                      },
                      name: { 
                        type: 'string',
                        example: 'John Doe'
                      },
                      email: {
                        type: 'string',
                        example: 'john@example.com'
                      },
                    },
                  },
                },
              },
            },
          },
        },
      },
      post: {
        summary: 'Create a new user',
        description: 'Create a new user in the system',
        tags: ['Users'],
        requestBody: {
          required: true,
          content: {
            'application/json': {
              schema: {
                type: 'object',
                required: ['name', 'email'],
                properties: {
                  name: {
                    type: 'string',
                    example: 'Jane Smith',
                  },
                  email: {
                    type: 'string',
                    format: 'email',
                    example: 'jane@example.com',
                  },
                  age: {
                    type: 'integer',
                    minimum: 0,
                    example: 28,
                  },
                },
              },
            },
          },
        },
        responses: {
          '201': {
            description: 'User created successfully',
            content: {
              'application/json': {
                schema: {
                  type: 'object',
                  properties: {
                    id: { type: 'integer' },
                    name: { type: 'string' },
                    email: { type: 'string' },
                    createdAt: { type: 'string', format: 'date-time' },
                  },
                },
              },
            },
          },
          '400': {
            description: 'Bad request - invalid input',
          },
        },
      },
    },
    '/api/v1/users/{id}': {
      get: {
        summary: 'Get user by ID',
        tags: ['Users'],
        parameters: [
          {
            name: 'id',
            in: 'path',
            required: true,
            schema: {
              type: 'integer',
            },
            description: 'User ID',
          },
        ],
        responses: {
          '200': {
            description: 'User found',
          },
          '404': {
            description: 'User not found',
          },
        },
      },
    },
  },
  components: {
    schemas: {
      User: {
        type: 'object',
        properties: {
          id: {
            type: 'integer',
            description: 'Unique identifier for the user',
          },
          name: {
            type: 'string',
            description: 'Full name of the user',
          },
          email: {
            type: 'string',
            format: 'email',
            description: 'Email address of the user',
          },
          createdAt: {
            type: 'string',
            format: 'date-time',
            description: 'When the user was created',
          },
        },
      },
    },
  },
};

// Middleware
app.use(express.json());

// Sample routes
app.get('/api/v1/users', (req, res) => {
  res.json([
    { id: 1, name: 'John Doe', email: 'john@example.com' },
    { id: 2, name: 'Jane Smith', email: 'jane@example.com' },
  ]);
});

app.post('/api/v1/users', (req, res) => {
  const { name, email, age } = req.body;
  const newUser = {
    id: Math.floor(Math.random() * 1000),
    name,
    email,
    age,
    createdAt: new Date().toISOString(),
  };
  res.status(201).json(newUser);
});

app.get('/api/v1/users/:id', (req, res) => {
  const { id } = req.params;
  res.json({ id: parseInt(id), name: 'Sample User', email: 'user@example.com' });
});

// Serve Scalar API Reference - CORRECT USAGE
app.use(
  '/reference',
  apiReference({
    spec: {
      content: openApiSpec,
    },
    theme: 'bluePlanet',
  })
);

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'OK', timestamp: new Date().toISOString() });
});

app.listen(PORT, () => {
  console.log(`🚀 Server running on http://localhost:${PORT}`);
  console.log(`📚 API Documentation: http://localhost:${PORT}/reference`);
});