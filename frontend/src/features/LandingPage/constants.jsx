import dayjs from 'dayjs';
import { CakeRounded, EmailRounded, LockRounded, SupervisedUserCircleRounded } from '@mui/icons-material';

/**
 * Sample Data that is used to populate the landing page graph
 */
export const SAMPLE_DATA = [
  { name: 'Entertainment', price: 320.0 },
  { name: 'Tools and equipment', price: 1451.99 },
  { name: 'Supplies', price: 758.99 },
  { name: 'Utilities', price: 420.0 },
  { name: 'Clothing', price: 430.0 },
  { name: 'Maintenance', price: 200.0 },
  { name: 'Education', price: 1341.59 },
  { name: 'Misc', price: 300.0 },
];

/**
 * Sample colors that is used to populate the landing page graph
 */
export const COLORS = [
  'rgba(255, 99, 132, 0.6)', // Red
  'rgba(54, 162, 235, 0.6)', // Blue
  'rgba(255, 206, 86, 0.6)', // Yellow
  'rgba(75, 192, 192, 0.6)', // Green
  'rgba(153, 102, 255, 0.6)', // Purple
  'rgba(255, 159, 64, 0.6)', // Orange
  'rgba(255, 99, 132, 0.6)', // Red
  'rgba(54, 162, 235, 0.6)', // Blue
  'rgba(255, 206, 86, 0.6)', // Yellow
];

/**
 * Modal State
 */
export const MODAL_STATE = {
  NONE: 'none',
  SIGN_UP: 'sign_up',
  SIGN_IN: 'sign_in',
  FORGOT_PASSWORD: 'forgot_password',
};

export const GENERIC_TEXTFIELD_VARIANT = {
  variant: 'outlined',
};

export const AUTOCOMPLETE_PASSWORD_PROPS = {
  autocomplete: 'true',
};

export const LOGIN_FORM_FIELDS = {
  email: {
    label: 'Email Address',
    placeholder: 'Enter your email address in correct format',
    value: '',
    name: 'email',
    size: 'small',
    errorMsg: '',
    required: true,
    type: 'email',
    icon: <EmailRounded />,
    fullWidth: true,
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Email is required',
      },
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Email should be less than 50 characters',
      },
      {
        validate: (value) => {
          var re = /\S+@\S+\.\S+/;
          return !re.test(value);
        },
        message: 'Email is not in the correct form',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
  },
  password: {
    label: 'Password',
    placeholder: 'Enter your secure password',
    value: '',
    type: 'password',
    name: 'password',
    size: 'small',
    errorMsg: '',
    icon: <LockRounded />,
    required: true,
    fullWidth: true,
    validators: [
      {
        validate: (value) => value.trim().length <= 6,
        message: 'Password is not strong enough.',
      },
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Password should be less than 50 characters',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
    ...AUTOCOMPLETE_PASSWORD_PROPS,
  },
};

export const SIGN_UP_FORM_FIELDS = {
  username: {
    label: 'Username',
    placeholder: 'Enter a username with at least four characters',
    value: '',
    name: 'username',
    size: 'small',
    errorMsg: '',
    required: true,
    type: 'text',
    icon: <SupervisedUserCircleRounded />,
    fullWidth: true,
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Username is required',
      },
      {
        validate: (value) => value.trim().length <= 3,
        message: 'Username cannot be less than 4 characters',
      },
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Username should be less than 50 characters',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
  },
  email: {
    label: 'Email Address',
    placeholder: 'Enter your unique email address in correct format',
    value: '',
    name: 'email',
    size: 'small',
    errorMsg: '',
    required: true,
    type: 'email',
    icon: <EmailRounded />,
    fullWidth: true,
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Email is required',
      },
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Email should be less than 50 characters',
      },
      {
        validate: (value) => {
          var re = /\S+@\S+\.\S+/;
          return !re.test(value);
        },
        message: 'Email is not in the correct form',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
  },
  password: {
    label: 'Password',
    placeholder: 'Enter your secure password with at least six characters',
    value: '',
    type: 'password',
    name: 'password',
    size: 'small',
    errorMsg: '',
    icon: <LockRounded />,
    required: true,
    fullWidth: true,
    validators: [
      {
        validate: (value) => value.trim().length <= 6,
        message: 'Password is not strong enough.',
      },
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Password should be less than 50 characters',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
    ...AUTOCOMPLETE_PASSWORD_PROPS,
  },
  birthday: {
    label: 'Date of Birth',
    placeholder: '',
    value: '',
    type: 'date',
    size: 'small',
    name: 'birthday',
    errorMsg: '',
    icon: <CakeRounded />,
    required: true,
    fullWidth: true,
    validators: [
      {
        validate: (birthdate) => {
          const currentYear = dayjs().year();
          const dob = dayjs(birthdate, 'YYYY-MM-DD');

          if (dob.isValid()) {
            const age = currentYear - dob.year() - dayjs().isBefore(dob, 'year');
            return age <= 13;
          }
          return false;
        },
        message: 'User must be above 13 years old to register.',
      },
      {
        validate: (value) => value.trim().length <= 2,
        message: 'User must add birthday.',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
  },
};

export const RESET_PASSWORD_FIELDS = {
  email: {
    label: 'Email Address',
    placeholder: 'Enter your email address',
    value: '',
    name: 'email',
    size: 'small',
    errorMsg: '',
    required: true,
    type: 'email',
    icon: <EmailRounded />,
    fullWidth: true,
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Email is required',
      },
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Email should be less than 50 characters',
      },
      {
        validate: (value) => {
          var re = /\S+@\S+\.\S+/;
          return !re.test(value);
        },
        message: 'Email is not in the correct form',
      },
    ],
    ...GENERIC_TEXTFIELD_VARIANT,
  },
};
