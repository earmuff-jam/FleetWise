import * as XLSX from 'xlsx';
import { useState } from 'react';
import { useDispatch } from 'react-redux';

import { enqueueSnackbar } from 'notistack';
import { Button, Stack } from '@mui/material';
import { SaveRounded } from '@mui/icons-material';

import { inventoryActions } from '@features/Assets/inventorySlice';
import ViewFileContent from '@features/Assets/AddAssetsInBulk/ViewFileContent';
import AddAssetsInBulkActions from '@features/Assets/AddAssetsInBulk/AddAssetsInBulkActions';
import { buildXcel } from '@common/utils';
import { BULK_ASSETS_HEADERS } from '@features/Assets/constants';

export default function AddAssetsInBulk({ handleClose }) {
  const dispatch = useDispatch();
  const [uploadedFileInJson, setUploadedFileInJson] = useState([]);
  const [fileDetails, setFileDetails] = useState({ name: '', lastModifiedDate: '', size: '' });

  const handleDownload = () => {
    const headers = BULK_ASSETS_HEADERS.map((v) => v.label);
    const templatedData = BULK_ASSETS_HEADERS.map((v) => v.value || '');
    buildXcel(headers, templatedData, 'asset-template.xlsx', 'assets');
  };

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setFileDetails({ name: file.name, lastModifiedDate: file.lastModifiedDate, size: file.size });
      const reader = new FileReader();
      reader.onload = (e) => {
        const data = e.target.result;
        const workbook = XLSX.read(data, { type: 'array' });
        const sheetName = workbook.SheetNames[0];
        const worksheet = workbook.Sheets[sheetName];
        const formattedArr = XLSX.utils.sheet_to_json(worksheet, {
          rawNumbers: true,
        });
        setUploadedFileInJson(formattedArr);
      };
      reader.readAsArrayBuffer(file);
    }
    // reset the target to null to support upload of the same file
    event.target.value = null;
  };

  const handleRemove = () => {
    setFileDetails({ name: '', lastModifiedDate: '', size: '' });
    setUploadedFileInJson(null);
  };

  const resetData = () => {
    setUploadedFileInJson(null);
    handleClose();
  };

  const submit = () => {
    if (Array.isArray(uploadedFileInJson) && uploadedFileInJson.length > 0) {
      dispatch(inventoryActions.addBulkInventory(uploadedFileInJson));
    }
    enqueueSnackbar('Uploaded inventories in bulk.', {
      variant: 'success',
    });
    resetData();
  };

  return (
    <Stack alignItems="center" spacing={1}>
      <AddAssetsInBulkActions
        handleFileChange={handleFileChange}
        fileDetails={fileDetails}
        handleClick={handleDownload}
      />
      <ViewFileContent
        handleRemove={handleRemove}
        showContent={Boolean(fileDetails?.name.length)}
        name={fileDetails.name}
        lastModifiedDate={fileDetails.lastModifiedDate}
        size={fileDetails.size}
      />
      <Button
        variant="outlined"
        startIcon={<SaveRounded />}
        onClick={submit}
        disabled={Boolean(!fileDetails?.name.length)}
      >
        Save
      </Button>
    </Stack>
  );
}
