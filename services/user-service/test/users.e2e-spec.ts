import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import request from 'supertest';
import { AppModule } from '../src/app.module';

describe('User Service E2E Tests', () => {
  let app: INestApplication;
  let userId: string;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    await app.init();
  });

  afterAll(async () => {
    await app.close();
  });

  describe('Health Check', () => {
    it('GET /health should return service status', () => {
      return request(app.getHttpServer())
        .get('/health')
        .expect(200)
        .expect((res) => {
          expect(res.body.status).toBe('healthy');
          expect(res.body.service).toBe('user-service');
        });
    });
  });

  describe('Users - Create', () => {
    it('POST /users should create a new user', () => {
      return request(app.getHttpServer())
        .post('/users')
        .send({
          email: 'john@example.com',
          firstName: 'John',
          lastName: 'Doe',
          password: 'password123',
        })
        .expect(201)
        .expect((res) => {
          expect(res.body.email).toBe('john@example.com');
          expect(res.body.firstName).toBe('John');
          expect(res.body.lastName).toBe('Doe');
          expect(res.body.id).toBeDefined();
          expect(res.body.createdAt).toBeDefined();
          userId = res.body.id;
        });
    });

    it('POST /users should fail with duplicate email', () => {
      return request(app.getHttpServer())
        .post('/users')
        .send({
          email: 'john@example.com',
          firstName: 'Jane',
          lastName: 'Doe',
          password: 'password456',
        })
        .expect(409);
    });
  });

  describe('Users - Read', () => {
    it('GET /users should return all users', () => {
      return request(app.getHttpServer())
        .get('/users')
        .expect(200)
        .expect((res) => {
          expect(Array.isArray(res.body)).toBe(true);
          expect(res.body.length).toBeGreaterThan(0);
        });
    });

    it('GET /users/:id should return a specific user', () => {
      return request(app.getHttpServer())
        .get(`/users/${userId}`)
        .expect(200)
        .expect((res) => {
          expect(res.body.id).toBe(userId);
          expect(res.body.email).toBe('john@example.com');
        });
    });

    it('GET /users/:id should return 404 for non-existent user', () => {
      return request(app.getHttpServer())
        .get('/users/non-existent-id')
        .expect(404);
    });
  });

  describe('Users - Update', () => {
    it('PATCH /users/:id should update user details', () => {
      return request(app.getHttpServer())
        .patch(`/users/${userId}`)
        .send({
          firstName: 'Jonathan',
          lastName: 'Smith',
        })
        .expect(200)
        .expect((res) => {
          expect(res.body.firstName).toBe('Jonathan');
          expect(res.body.lastName).toBe('Smith');
          expect(res.body.email).toBe('john@example.com');
        });
    });
  });

  describe('Users - Delete', () => {
    it('DELETE /users/:id should remove a user', () => {
      return request(app.getHttpServer())
        .delete(`/users/${userId}`)
        .expect(200)
        .expect((res) => {
          expect(res.body.id).toBe(userId);
        });
    });

    it('GET /users/:id should return 404 after deletion', () => {
      return request(app.getHttpServer()).get(`/users/${userId}`).expect(404);
    });
  });
});
