import { ValidationErrors } from '@angular/forms';
import { ValidationError } from 'class-validator';

export const errorsToFormErrors = (
  errors: ValidationError[]
): ValidationErrors => {
  return errors.reduce(
    (res, { property, constraints }) => ({
      ...res,
      [property]: constraints ? Object.values(constraints)[0] : null,
    }),
    {} as ValidationErrors
  );
};
