import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
} from '@nestjs/common';
import { UsersService } from './users.service';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { UserResponseDto } from './dto/user-response.dto';

@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Post()
  create(@Body() createUserDto: CreateUserDto): UserResponseDto {
    const user = this.usersService.create(createUserDto);
    return this.formatUserResponse(user);
  }

  @Get()
  findAll(): UserResponseDto[] {
    const users = this.usersService.findAll();
    return users.map((user) => this.formatUserResponse(user));
  }

  @Get(':id')
  findOne(@Param('id') id: string): UserResponseDto {
    const user = this.usersService.findOne(id);
    return this.formatUserResponse(user);
  }

  @Patch(':id')
  update(
    @Param('id') id: string,
    @Body() updateUserDto: UpdateUserDto,
  ): UserResponseDto {
    const user = this.usersService.update(id, updateUserDto);
    return this.formatUserResponse(user);
  }

  @Delete(':id')
  remove(@Param('id') id: string): UserResponseDto {
    const user = this.usersService.remove(id);
    return this.formatUserResponse(user);
  }

  private formatUserResponse(user: any): UserResponseDto {
    return {
      id: user.id,
      email: user.email,
      firstName: user.firstName,
      lastName: user.lastName,
      createdAt: user.createdAt,
      updatedAt: user.updatedAt,
    };
  }
}
