import React, { useEffect, useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Button, TextField, Tooltip, Typography } from '@material-ui/core';

import { produce } from 'immer';
import { enqueueSnackbar } from 'notistack';
import { ADD_EXPENSE_FORM_FIELDS } from './constants';
import { useDispatch, useSelector } from 'react-redux';
import { eventActions } from '../../Containers/Event/eventSlice';
import Autocomplete, { createFilterOptions } from '@material-ui/lab/Autocomplete';

const useStyles = makeStyles((theme) => ({
  container: {
    display: 'flex',
    flexDirection: 'column',
    gap: theme.spacing(6),
    padding: theme.spacing(4),
  },
}));

const filter = createFilterOptions();

const AddExpenseDetail = ({ eventID, userID, setDisplayMode }) => {
  const classes = useStyles();
  const dispatch = useDispatch();

  const { loading, categories } = useSelector((state) => state.event);

  const [options, setOptions] = useState([]);
  const [selectedCategoryName, setSelectedCategoryName] = useState(null);
  const [selectedCategory, setSelectedCategory] = useState({});
  const [formFields, setFormFields] = useState(ADD_EXPENSE_FORM_FIELDS);

  const handleInput = (event) => {
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

  const handleFormSubmit = (e) => {
    e.preventDefault();

    const containsErr = Object.values(formFields).reduce((acc, el) => {
      if (el.errorMsg) {
        return true;
      }
      return acc;
    }, false);

    const requiredFormFields = Object.values(formFields).filter((v) => v.required);
    const isRequiredFieldsEmpty = requiredFormFields.some((el) => el.value.trim() === '');

    if (
      containsErr ||
      isRequiredFieldsEmpty ||
      selectedCategory === null ||
      selectedCategoryName === null ||
      Object.keys(selectedCategory).length <= 0
    ) {
      console.log('Empty form fields. Unable to proceed.');
      enqueueSnackbar('Unable to add new expense report.', {
        variant: 'error',
      });
      return;
    } else {
      const formattedData = Object.values(formFields).reduce((acc, el) => {
        if (el.value) {
          acc[el.name] = el.value;
        }
        return acc;
      }, {});

      const categoryNameOrID =
        categories.find((v) => v.category_name === selectedCategory.selectedCategory)?.id ||
        selectedCategory.selectedCategory;

      const postFormattedData = {
        ...formattedData,
        category_id: categoryNameOrID,
        category_name: selectedCategoryName,
        created_by: userID,
        eventID: eventID,
      };

      dispatch(eventActions.addExpenseList({ postFormattedData }));
      setFormFields(ADD_EXPENSE_FORM_FIELDS);
      setDisplayMode();
      enqueueSnackbar('Successfully added new expense report.', {
        variant: 'success',
      });
    }
  };

  useEffect(() => {
    if (!loading && Array.isArray(categories)) {
      setOptions(categories.map((v) => ({ category_id: v.id, category_name: v.category_name })));
    }
  }, [categories, loading]);

  return (
    <div className={classes.container}>
      {Object.values(formFields).map((v) => (
        <TextField
          key={v.name}
          id={v.name}
          className={classes.text}
          variant={v.variant}
          fullWidth={v.fullWidth}
          name={v.name}
          label={v.label}
          required={v.required}
          value={v.value}
          placeholder={v.placeholder}
          onChange={handleInput}
          error={!!v.errorMsg}
          helperText={v.errorMsg}
          multiline={v.multiline || false}
          rows={v.rows || 4}
        />
      ))}
      <Autocomplete
        value={selectedCategory.inputValue}
        forcePopupIcon
        onChange={(event, newValue) => {
          if (typeof newValue === 'string') {
            setSelectedCategory({
              selectedCategory: newValue,
            });
            setSelectedCategoryName(newValue);
          } else if (newValue && newValue.inputValue) {
            setSelectedCategory({
              selectedCategory: newValue.inputValue,
            });
            setSelectedCategoryName(newValue.inputValue);
          }
        }}
        filterOptions={(options, params) => {
          const filtered = filter(options, params);
          if (params.inputValue !== '') {
            filtered.push({
              inputValue: params.inputValue,
              category: `Add "${params.inputValue}"`,
            });
          }

          return filtered;
        }}
        selectOnFocus
        clearOnBlur
        handleHomeEndKeys
        id="autocomplete-category-details"
        options={options.map((v) => v.category_name)}
        getOptionLabel={(option) => (typeof option === 'string' ? option : option.category || '')}
        renderOption={(option) => option.category || option}
        freeSolo
        renderInput={(params) => (
          <TextField {...params} fullWidth label="Category" variant="standard" placeholder="Categorized as ..." />
        )}
      />
      <Tooltip
        title={
          'Add expense report to the selected event. All expense reports must be approved by the creator of the selected event.'
        }
      >
        <Button onClick={handleFormSubmit}>Submit</Button>
      </Tooltip>
    </div>
  );
};

export default AddExpenseDetail;
