import { Injectable } from '@nestjs/common';

@Injectable()
export class WeatherService {
  
  // Função que o controller está procurando para CRIAR
  create(createWeatherDto: any) {
    return 'Essa ação adiciona um novo clima (mock)';
  }

  // Função que o controller está procurando para BUSCAR TUDO
  findAll() {
    return 'Essa ação retorna todos os climas (mock)';
  }
}