import ForgotPasswordText from '@features/LandingPage/Authentication/Login/ForgotPasswordText';

export default {
  title: 'LandingPage/Authentication/Login/ForgotPasswordText',
  component: ForgotPasswordText,
  tags: ['autodocs'],
};

const Template = (args) => <ForgotPasswordText {...args} />;

export const ForgotPasswordTextDefault = Template.bind({});

ForgotPasswordTextDefault.args = {
  handleForgotPassword: () => {},
};
