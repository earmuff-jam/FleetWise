import { useEffect, useState } from 'react';
import { enqueueSnackbar } from 'notistack';
import { FORGOT_PASSWORD_FORM_FIELDS } from '../../constants';
import { Button, List, ListItem, ListItemText, Skeleton, Stack, TextField, Typography } from '@mui/material';
import { useSelector } from 'react-redux';
import RowHeader from '@common/RowHeader';

export const SECURITY_QUESTIONS = [
  {
    id: 1,
    label: 'Which year did you graduate high school ?',
  },
  {
    id: 2,
    label: 'Where did your parents first meet ?',
  },
];

export default function SecurityQuestionAnswerForm({ handleCloseModal }) {
  const { isUserVerified, retryCounts, loading } = useSelector((state) => state.auth);

  const [verifiedUser, setVerifiedUser] = useState(false);
  const [formFields, setFormFields] = useState(FORGOT_PASSWORD_FORM_FIELDS);

  const handleInput = (event) => {
    const { name, value } = event.target;
    const updatedFormFields = Object.assign({}, formFields, {
      [name]: {
        ...formFields[name],
        value: value,
        errorMsg: '',
      },
    });

    for (const validator of updatedFormFields[name].validators) {
      if (validator.validate(value)) {
        updatedFormFields[name].errorMsg = validator.message;
        break;
      }
    }
    setFormFields(updatedFormFields);
  };

  const resetData = () => {
    setFormFields({ ...FORGOT_PASSWORD_FORM_FIELDS });
    handleCloseModal();
  };

  const handleSubmit = () => {
    const containsErr = Object.values(formFields).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formFields).filter((v) => v.required);
    const isRequiredFieldsEmpty = requiredFormFields.some((el) => el.value.trim() === '');

    if (containsErr || isRequiredFieldsEmpty) {
      enqueueSnackbar('Cannot continue.', {
        variant: 'error',
      });
      return;
    }

    // const formattedData = Object.values(formFields).reduce((acc, el) => {
    //   if (el.value) {
    //     acc[el.name] = el.value;
    //   }
    //   return acc;
    // }, {});

    // const draftRequest = {
    //   ...formattedData,
    //   created_on: dayjs().toISOString(),
    // };
    //   dispatch(maintenancePlanActions.createPlan(draftRequest));
    resetData();
  };

  useEffect(() => {
    if (!loading && isUserVerified) {
      setVerifiedUser(true);
    }
    return () => {
      setVerifiedUser(false);
    };
  }, [loading]);

  if (loading) {
    return <Skeleton height="10vh" />;
  }

  if (verifiedUser) {
    return (
      <>
        <RowHeader
          title="Security questions"
          caption="Please answer your security questions to the best of your knowledge."
        />
        <List>
          <ListItem key={SECURITY_QUESTIONS[0].id}>
            <Stack width="inherit">
              <ListItemText primary={<Typography>{SECURITY_QUESTIONS[0].label}</Typography>} />
              <TextField
                key={formFields['question_01'].name}
                id={formFields['question_01'].name}
                name={formFields['question_01'].name}
                label={formFields['question_01'].label}
                value={formFields['question_01'].value}
                placeholder={formFields['question_01'].placeholder}
                onChange={handleInput}
                required={formFields['question_01'].required}
                fullWidth={formFields['question_01'].fullWidth}
                error={!!formFields['question_01'].errorMsg}
                helperText={formFields['question_01'].errorMsg}
                variant={formFields['question_01'].variant}
                minRows={formFields['question_01'].rows || 4}
                multiline={formFields['question_01'].multiline || false}
              />
            </Stack>
          </ListItem>
          <ListItem key={SECURITY_QUESTIONS[1].id}>
            <Stack width="inherit">
              <ListItemText primary={<Typography>{SECURITY_QUESTIONS[1].label}</Typography>} />
              <TextField
                key={formFields['question_02'].name}
                id={formFields['question_02'].name}
                name={formFields['question_02'].name}
                label={formFields['question_02'].label}
                value={formFields['question_02'].value}
                placeholder={formFields['question_02'].placeholder}
                onChange={handleInput}
                required={formFields['question_02'].required}
                fullWidth={formFields['question_02'].fullWidth}
                error={!!formFields['question_02'].errorMsg}
                helperText={formFields['question_02'].errorMsg}
                variant={formFields['question_02'].variant}
                minRows={formFields['question_02'].rows || 4}
                multiline={formFields['question_02'].multiline || false}
              />
            </Stack>
          </ListItem>
          <Button variant="outlined" onClick={handleSubmit}>
            Submit
          </Button>
        </List>
      </>
    );
  } else if (retryCounts >= 3) {
    return <Typography color="error.main"> User cannot be verified. Remaining attempts {3 - retryCounts}</Typography>;
  } else {
    return <Typography color="error.main"> Unable to verify user. Max password reset exceeded.</Typography>;
  }
}
