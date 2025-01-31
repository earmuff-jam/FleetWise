import { RESET_PASSWORD_FIELDS } from '@features/LandingPage/constants';
import ForgotPasswordFormFields from '@features/LandingPage/Authentication/ForgotPassword/ForgotPasswordFormFields';

export default {
  title: 'LandingPage/Authentication/ForgotPassword/ForgotPasswordFormFields',
  component: ForgotPasswordFormFields,
  tags: ['autodocs'],
};

const Template = (args) => <ForgotPasswordFormFields {...args} />;

export const ForgotPasswordFormFieldsDefault = Template.bind({});

ForgotPasswordFormFieldsDefault.args = {
  formFields: RESET_PASSWORD_FIELDS,
  handleInput: () => {},
};
