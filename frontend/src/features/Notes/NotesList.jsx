import { useEffect, useState } from 'react';

import { useDispatch, useSelector } from 'react-redux';

import { Stack } from '@mui/material';
import Note from '@features/Notes/Note/Note';
import SimpleModal from '@common/SimpleModal';
import AddNote from '@features/Notes/AddNote/AddNote';
import { notesActions } from '@features/Notes/notesSlice';
import NoteHeader from '@features/Notes/Header/NoteHeader';

const NotesList = () => {
  const dispatch = useDispatch();
  const { loading, notes } = useSelector((state) => state.notes);

  const [editMode, setEditMode] = useState(false);
  const [selecteNoteID, setSelectedNoteID] = useState(null);

  const handleEditMode = () => setEditMode(!editMode);

  const resetData = () => {
    setEditMode(false);
    setSelectedNoteID(null);
  };

  useEffect(() => {
    dispatch(notesActions.getNotes());
  }, []);

  return (
    <Stack spacing="1rem" data-tour="notes-0">
      <NoteHeader handleClick={handleEditMode} />
      <Note notes={notes} loading={loading} setEditMode={setEditMode} setSelectedNoteID={setSelectedNoteID} />
      {editMode && (
        <SimpleModal
          title="Add New Note"
          subtitle={'Notes with assigned locations reflect approximate values.'}
          handleClose={resetData}
          maxSize="xs"
        >
          <AddNote
            setEditMode={setEditMode}
            setSelectedNoteID={setSelectedNoteID}
            noteID={selecteNoteID}
            notes={notes}
          />
        </SimpleModal>
      )}
    </Stack>
  );
};

export default NotesList;
