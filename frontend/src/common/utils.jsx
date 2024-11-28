import { Button, Stack, Typography } from '@mui/material';
import SimpleModal from './SimpleModal';
import CryptoJS from 'crypto-js';

/**
 * Function used to decrypt the client location from the local storage. To prevent XSS attacks on
 * the client web provider, we are using CyrptoJS lib to persist a scrambled version of the user location.
 *
 * @returns parsed and validated user location data
 */
export const RetrieveClientLocation = () => {
  const ciphertext = localStorage.getItem('client_location');
  if (ciphertext) {
    try {
      const userID = localStorage.getItem('userID');
      const bytes = CryptoJS.AES.decrypt(ciphertext, userID);
      const decryptedData = JSON.parse(bytes.toString(CryptoJS.enc.Utf8));
      return decryptedData || {};
    } catch (e) {
      // eat exception
    }
  }
  return null;
};

/**
 * display no matching records found component if there are no records
 */
export const EmptyComponent = ({ subtitle = '' }) => (
  <Stack alignItems="center">
    <Typography color="textSecondary">Sorry, no matching records found.</Typography>
    <Typography variant="caption" color="textSecondary">
      {subtitle}
    </Typography>
  </Stack>
);

/**
 * Confirmation Box Modal
 */
export const ConfirmationBoxModal = ({ openDialog, title, handleClose, maxSize, deleteID, confirmDelete }) => {
  return openDialog ? (
    <SimpleModal title={title} handleClose={handleClose} maxSize={maxSize}>
      <Typography variant="subtitle2">Delete this item?</Typography>
      <Stack direction="row" justifyContent="flex-end" spacing={1}>
        <Button onClick={handleClose}>Go back</Button>
        <Button variant="outlined" autoFocus onClick={() => confirmDelete(deleteID)}>
          Confirm
        </Button>
      </Stack>
    </SimpleModal>
  ) : null;
};

/**
 * Takes array of notes and transforms them into objects categorized by the date and time
 *
 * Recently updated - Up to the last week
 * Last Week - Up to the previous week
 * Last month and beyond - everything else
 *
 * @param {Array} notes
 * @returns {Array} refactored notes
 */
export const categorizeNotes = (notes) => {
  const currentTime = new Date();
  const categorizedNotes = notes.reduce((acc, item) => {
    const updatedTime = new Date(item.updated_at);
    const differenceInDays = Math.floor((currentTime - updatedTime) / (1000 * 3600 * 24));

    let category;
    if (differenceInDays <= 7) {
      category = 'Recently added notes';
    } else if (differenceInDays <= 14) {
      category = 'Last Week';
    } else {
      category = 'Last Month and Beyond';
    }

    if (!acc[category]) {
      acc[category] = {
        id: acc.length + 1,
        category: category,
        totalNotes: 0,
        details: [],
      };
    }

    acc[category].details.push({
      id: acc[category].details.length + 1,
      ...item,
    });

    acc[category].totalNotes++;
    return acc;
  }, {});

  return Object.values(categorizedNotes);
};

/**
 * Function used to capitalize the first letter of passed in string
 * @param {word} string - the word to capitalize
 * @returns string - with captial first letter
 */
export const capitalizeFirstLetter = (string) => {
  return string.charAt(0).toUpperCase() + string.slice(1);
};

/**
 * Function returns a string with prefix if prefix exists.
 * @param {string} prefixVal - the prefix value to append
 * @param {any} value - the value to append the prefix string onto
 * @param {roundUpTo} int - rounds to certain number if the selected value is a number
 * @returns any - result of combination of prefixValue and value
 */
export const prefix = (prefixVal, value, roundUpTo = 2) => {
  if (prefixVal != null) {
    if (Number(value)) {
      return `${prefixVal} ${Number(value).toFixed(roundUpTo)}`;
    } else {
      return `${prefixVal} ${value}`;
    }
  }
  return value;
};

/**
 * Function used to pluralize the selected word. Eg, Dog -> dogs.
 * @param {string} stringToEdit - The string that is to be pluralized
 * @param {int} size - The number that determines if the words needs to be pluralized or not
 * @returns {string} - The new string that is pluralized if necessary
 */
export const pluralizeWord = (stringToEdit, size) => {
  if (size <= 1) {
    return `${size} ${stringToEdit}`;
  } else {
    return `${size} ${stringToEdit}s`;
  }
};
