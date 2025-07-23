import { Expose } from 'class-transformer';
import {  IsNotEmpty,  MaxLength, MinLength } from 'class-validator';

export class Auth {
  @MinLength(6)
  @MaxLength(256)
  @Expose()
  login!: string;

  @MinLength(6)
  @MaxLength(256)
  @IsNotEmpty()
  @Expose()
  password!: string;
}
