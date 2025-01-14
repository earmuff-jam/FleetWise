import { useEffect, useState } from 'react';

import dayjs from 'dayjs';
import { useDispatch } from 'react-redux';
import { enqueueSnackbar } from 'notistack';

import { Accordion, Skeleton } from '@mui/material';
import relativeTime from 'dayjs/plugin/relativeTime';
import { notesActions } from '@features/Notes/notesSlice';
import { categorizeNotes, ConfirmationBoxModal, EmptyComponent } from '@common/utils';
import NoteAccordionSummary from '@features/Notes/NoteAccordion/NoteAccordionSummary';
import NoteAccordionDetails from '@features/Notes/NoteAccordion/NoteAccordionDetails';

const Note = ({ notes, loading, setEditMode, setSelectedNoteID }) => {
  const dispatch = useDispatch();
  dayjs.extend(relativeTime);

  const [deleteID, setDeleteID] = useState(-1);
  const [formattedNotes, setFormattedNotes] = useState([]);
  const [confirmDelete, setConfirmDelete] = useState(false);

  const reset = () => {
    setConfirmDelete(false);
    setDeleteID(-1);
  };

  const handleConfirmDelete = () => {
    removeSelectedNote(deleteID);
    setConfirmDelete(false);
  };

  const removeSelectedNote = (noteID) => {
    dispatch(notesActions.removeNote({ noteID }));
    enqueueSnackbar('Successfully removed notes.', {
      variant: 'success',
    });
  };

  useEffect(() => {
    if (Array.isArray(notes) && notes.length >= 0) {
      const transformedData = categorizeNotes(notes);
      setFormattedNotes(transformedData);
    }
  }, [loading, notes]);

  if (loading) {
    return <Skeleton height="2rem" />;
  }

  if (!notes || notes.length === 0) {
    return <EmptyComponent subtitle="Add new notes." />;
  }

  return (
    <>
      {formattedNotes.map((v, index) => (
        <Accordion key={index} elevation={0} defaultExpanded={index === 0 ? true : false}>
          <NoteAccordionSummary title={v.category} totalNotes={v.totalNotes} color={v.color} />
          <NoteAccordionDetails
            details={v.details}
            setEditMode={setEditMode}
            setSelectedNoteID={setSelectedNoteID}
            setConfirmDelete={setConfirmDelete}
            setDeleteID={setDeleteID}
          />
        </Accordion>
      ))}
      <ConfirmationBoxModal
        openDialog={confirmDelete}
        title="Confirm deletion"
        handleClose={reset}
        maxSize="xs"
        deleteID={deleteID}
        confirmDelete={handleConfirmDelete}
      />
    </>
  );
};

export default Note;
