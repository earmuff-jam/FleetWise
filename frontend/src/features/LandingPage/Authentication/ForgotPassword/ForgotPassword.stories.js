import ForgotPassword from '@features/LandingPage/Authentication/ForgotPassword/ForgotPassword';

export default {
  title: 'LandingPage/Authentication/ForgotPassword/ForgotPassword',
  component: ForgotPassword,
  tags: ['autodocs'],
};

const Template = (args) => <ForgotPassword {...args} />;

export const ForgotPasswordDefault = Template.bind({});

ForgotPasswordDefault.args = {
  handleClose: () => {},
};
