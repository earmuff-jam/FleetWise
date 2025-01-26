import { DownloadRounded } from '@mui/icons-material';
import { Button, Stack, Typography } from '@mui/material';

export default function AddAssetsInBulkActions({ handleFileChange, fileDetails, handleClick }) {
  return (
    <>
      <Typography variant="caption" alignSelf="flex-start">
        Use the provided template to ensure that all the fields are properly recorded into the system. Edit selected
        asset to add more information.
      </Typography>
      <Stack direction="row" spacing={1}>
        <Button
          variant="outlined"
          component="label"
          onClick={handleClick}
          startIcon={<DownloadRounded color="primary" />}
        >
          Download Template
        </Button>
        <Button
          variant="outlined"
          component="label"
          onChange={handleFileChange}
          disabled={Boolean(fileDetails?.name.length)}
        >
          Upload File
          <input type="file" hidden />
        </Button>
      </Stack>
    </>
  );
}
