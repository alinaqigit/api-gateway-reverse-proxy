import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getHealth(): object {
    return {
      status: 'healthy',
      service: 'user-service',
      timestamp: new Date().toISOString(),
    };
  }
}
