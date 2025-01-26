import dayjs from 'dayjs';
import * as XLSX from 'xlsx';
import CryptoJS from 'crypto-js';

import relativeTime from 'dayjs/plugin/relativeTime';
import { Button, Stack, Typography } from '@mui/material';

import SimpleModal from '@common/SimpleModal';

dayjs.extend(relativeTime);

/**
 * Function used to convert a valid base64 encoded string to a binary data set. Used
 * in the context of images that need to be displayed in the list of categories or list
 * of maintenance plans.
 *
 * @type {string} base64EncodedString - The string that is a valid representation of image.
 * @returns valid bytes of images in Uint8Array format.
 */
export const Base64ToUint8Array = (base64EncodedString) => {
  const binaryString = atob(base64EncodedString);
  const len = binaryString.length;
  const bytes = new Uint8Array(len);

  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  return bytes;
};

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
export const EmptyComponent = ({ subtitle = '', padding = '0rem 0rem' }) => (
  <Stack alignItems="center" sx={{ padding: padding }}>
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
        <Button variant="outlined" autoFocus onClick={() => confirmDelete(deleteID || true)}>
          Confirm
        </Button>
      </Stack>
    </SimpleModal>
  ) : null;
};

/**
 * Takes array of notes and transforms them into objects categorized by the date and time
 *
 * Recently Edited - Up to the last week
 * Edited within last week - Up to the previous week
 * Edited couple of months ago - everything else
 *
 * Color Primary is used for Recently Edited
 * Color Secondary is used for Edited within last week
 * Color Default is used for Edit couple of month ago
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
    let color;
    if (differenceInDays <= 7) {
      category = 'Recently Edited';
      color = 'primary';
    } else if (differenceInDays <= 14) {
      category = 'Edited within last week';
      color = 'secondary';
    } else {
      category = 'Edited couple of months ago';
      color = 'default';
    }

    if (!acc[category]) {
      acc[category] = {
        id: acc.length + 1,
        category: category,
        totalNotes: 0,
        color: color,
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
 * Function used to pluralize the selected word. If the word does not exist, it will
 * ignore the word. Eg, ('Dog', 2) -> 2 dogs, ('', 2) -> 2
 * @param {string} stringToEdit - The string that is to be pluralized
 * @param {int} size - The number that determines if the words needs to be pluralized or not
 * @returns {string} - The new string that is pluralized if necessary
 */
export const pluralizeWord = (stringToEdit = '', size) => {
  if (!stringToEdit) {
    return `${size}`;
  }
  if (size <= 1) {
    return `${size} ${stringToEdit}`;
  } else {
    return `${size} ${stringToEdit}s`;
  }
};

/**
 * Function is used to format the selected date. If the date does not exist, returns null.
 * @param {string} Date object valid for dayjs to parse.
 * @param {string} Layout layout of the date that the date object should be subjected to
 * @returns {string} - The new formatted date that is parsed with the given layout
 */
export const formatDate = (date, layout = 'MMM, YYYY') => {
  if (!date) return null;
  return dayjs(date).isValid() && `${dayjs(date).format(layout)}`;
};

/**
 * Function that is used to format the given file size into human readable format
 * @param {string} sizeInBytes - the actual size of the file in bytes format
 * @returns - The formatted size of the file size in human readable format.
 */
export function formatFileSize(sizeInBytes) {
  if (sizeInBytes === 0) return '0 Bytes';
  const units = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(sizeInBytes) / Math.log(1024));
  const size = (sizeInBytes / Math.pow(1024, i)).toFixed(2); // Keeps 2 decimal places
  return `${size} ${units[i]}`;
}

/**
 * Function that is used to build an xcel sheet. Returns the built excel
 * sheet. Can be used to populate CSV data or build template data for
 * adding bulk assets.
 *
 * @param {Array} headers - the array of headers for the title
 * @param {Array} data - the data to populate the selected sheet with
 * @param {string} fileName - the name of the file
 * @param {string} sheetName - the name of the sheet
 */
export const buildXcel = (headers = [], data = [], fileName, sheetName) => {
  const wb = XLSX.utils.book_new();

  const colWidths = headers.map((header) => ({
    wch: header.length + 2,
  }));

  const ws = XLSX.utils.json_to_sheet(data, { origin: 'A2', skipHeader: true });

  ws['!cols'] = colWidths;
  XLSX.utils.sheet_add_aoa(ws, [headers], { origin: 'A1' });
  XLSX.utils.book_append_sheet(wb, ws, sheetName);

  XLSX.writeFile(wb, fileName);
};
