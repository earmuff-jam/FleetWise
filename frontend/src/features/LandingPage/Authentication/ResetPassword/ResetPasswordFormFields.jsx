import { useState } from 'react';

import { Stack } from '@mui/material';

import { RESET_PASSWORD_FORM_FIELDS } from '@features/LandingPage/constants';
import TextFieldWithLabel from '@common/TextFieldWithLabel/TextFieldWithLabel';
import { produce } from 'immer';

export default function ResetPasswordFormFields() {
  const [formFields, setFormFields] = useState(RESET_PASSWORD_FORM_FIELDS);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormFields(
      produce(formFields, (draft) => {
        draft[name].value = value;
        draft[name].errorMsg = '';

        for (const validator of draft[name].validators) {
          if (validator.validate(value)) {
            draft[name].errorMsg = validator.message;
            break;
          }
        }
      })
    );
  };

  const validate = (formFields) => {
    const containsErr = Object.values(formFields).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formFields).filter((v) => v.required);
    return containsErr || requiredFormFields.some((el) => el.value.trim() === '');
  };

  const submit = (e) => {
    e.preventDefault();

    if (validate(formFields)) {
      return;
    } else {
      const formattedData = Object.values(formFields).reduce((acc, el) => {
        if (el.value) {
          acc[el.name] = el.value;
        }
        return acc;
      }, {});
      // dispatch(authActions.getUserID(formattedData));
    }
  };

  return (
    <Stack spacing={1} width={{ xs: 'calc(100% - 1rem)', sm: '50%' }}>
      {/* OTP generator */}
      <TextFieldWithLabel
        id={formFields.recovery_token.name}
        name={formFields.recovery_token.name}
        label={formFields.recovery_token.label}
        value={formFields.recovery_token.value}
        size={formFields.recovery_token.size}
        placeholder={formFields.recovery_token.placeholder}
        handleChange={handleInputChange}
        required={formFields.recovery_token.required}
        fullWidth={formFields.recovery_token.fullWidth}
        error={Boolean(formFields.recovery_token.errorMsg)}
        helperText={formFields.recovery_token.errorMsg}
        variant={formFields.recovery_token.variant}
      />

      {/* New password */}
      <TextFieldWithLabel
        id={formFields.password.name}
        name={formFields.password.name}
        label={formFields.password.label}
        value={formFields.password.value}
        size={formFields.password.size}
        placeholder={formFields.password.placeholder}
        handleChange={handleInputChange}
        required={formFields.password.required}
        fullWidth={formFields.password.fullWidth}
        error={Boolean(formFields.password.errorMsg)}
        helperText={formFields.password.errorMsg}
        variant={formFields.password.variant}
      />

      {/* Confirm Password */}
      <TextFieldWithLabel
        id={formFields.confirmPassword.name}
        name={formFields.confirmPassword.name}
        label={formFields.confirmPassword.label}
        value={formFields.confirmPassword.value}
        size={formFields.confirmPassword.size}
        placeholder={formFields.confirmPassword.placeholder}
        handleChange={handleInputChange}
        required={formFields.confirmPassword.required}
        fullWidth={formFields.confirmPassword.fullWidth}
        error={Boolean(formFields.confirmPassword.errorMsg)}
        helperText={formFields.confirmPassword.errorMsg}
        variant={formFields.confirmPassword.variant}
      />
    </Stack>
  );
}
