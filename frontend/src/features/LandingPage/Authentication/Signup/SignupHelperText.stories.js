import SignupHelperText from '@features/LandingPage/Authentication/Signup/SignupHelperText';

export default {
  title: 'LandingPage/Authentication/Signup/SignupHelperText',
  component: SignupHelperText,
  tags: ['autodocs'],
};

const Template = (args) => <SignupHelperText {...args} />;

export const SignupHelperTextDefault = Template.bind({});
export const SignupHelperTextLoading = Template.bind({});
export const SignupHelperTextUniqueEmailAddress = Template.bind({});

SignupHelperTextDefault.args = {
  isEmailUnique: false,
  loading: false,
};

SignupHelperTextLoading.args = {
  isEmailUnique: false,
  loading: true,
};

SignupHelperTextUniqueEmailAddress.args = {
  isEmailUnique: true,
  loading: false,
};
