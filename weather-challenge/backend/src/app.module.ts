import { Module } from '@nestjs/common';
import { WeatherModule } from './weather/weather.module';
import { UsersModule } from './users/users.module';
import { MongooseModule } from '@nestjs/mongoose';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    MongooseModule.forRoot(process.env.MONGODB_URI || 'mongodb://localhost/weather-challenge'),
    WeatherModule,
    UsersModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}