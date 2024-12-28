import { profileActions } from '@features/Profile/profileSlice';
import { Autocomplete, Button, Chip, Stack, TextField } from '@mui/material';
import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

export default function SharableGroups({ handleSubmit, existingGroups, creator }) {
  const dispatch = useDispatch();
  const { loading, profiles = [] } = useSelector((state) => state.profile);

  const [options, setOptions] = useState([]);
  const [sharableGroups, setSharableGroups] = useState([]);

  useEffect(() => {
    if (!loading && Array.isArray(profiles)) {
      const draftProfiles = profiles.map((v) => ({
        display: v.email_address,
        value: v.id,
        label: v.email_address,
      }));
      setOptions(draftProfiles);
    }
  }, [loading, profiles]);

  useEffect(() => {
    if (Array.isArray(existingGroups) && existingGroups.length > 0 && Array.isArray(profiles) && profiles.length > 0) {
      const profileMap = new Map(profiles.map((profile) => [profile.id, profile]));
      const collaborators = existingGroups
        ?.map((userID) => profileMap.get(userID))
        .filter((profile) => profile !== undefined);

      const draftCollaborators = collaborators.map((profile) => ({
        display: profile.email_address,
        value: profile.id,
        label: profile.email_address,
      }));

      setSharableGroups(draftCollaborators);
    }
  }, [existingGroups, profiles]);

  useEffect(() => {
    dispatch(profileActions.getProfileList());
  }, []);

  return (
    <Stack spacing="0.2rem">
      <Autocomplete
        id="sharable-groups-options"
        multiple
        freeSolo
        limitTags={5}
        loading={loading}
        options={options}
        value={sharableGroups}
        isOptionEqualToValue={(option, value) => {
          return option.value === value.value;
        }}
        onChange={(_, newValue) => {
          setSharableGroups(newValue);
        }}
        renderInput={(params) => (
          <TextField {...params} variant="standard" label="Add collaborators" placeholder="Collaborators" />
        )}
        renderTags={(tagValue, getTagProps) =>
          tagValue.map((option, index) => {
            const { key, ...tagProps } = getTagProps({ index });
            return <Chip key={key} label={option.display} {...tagProps} disabled={option.value === creator} />;
          })
        }
      />
      <Button variant="text" onClick={() => handleSubmit(sharableGroups)} disabled={sharableGroups.length === 0}>
        Submit
      </Button>
    </Stack>
  );
}
