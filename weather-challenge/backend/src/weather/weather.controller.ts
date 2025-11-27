import { Body, Controller, Get, Post } from '@nestjs/common';
import { WeatherService } from './weather.service';

@Controller('api/weather') // Prefixo da rota
export class WeatherController {
  constructor(private readonly weatherService: WeatherService) {}

  // Endpoint que o Worker Go vai chamar para salvar dados
  // POST http://localhost:3000/api/weather/logs
  @Post('logs')
  async createLog(@Body() body: any) {
    console.log('Recebido log do Worker Go:', body);
    return this.weatherService.create(body);
  }

  // Endpoint que o Frontend vai chamar para mostrar o dashboard
  // GET http://localhost:3000/api/weather/logs
  @Get('logs')
  async getLogs() {
    return this.weatherService.findAll();
  }
}