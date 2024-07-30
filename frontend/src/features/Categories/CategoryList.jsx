import { Stack } from '@mui/material';
import { useState } from 'react';
import { AddRounded } from '@mui/icons-material';
import SimpleModal from '../common/SimpleModal';
import AddCategory from './AddCategory';
import HeaderWithButton from '../common/HeaderWithButton';
import Category from './Category';
import { useSelector } from 'react-redux';

const CategoryList = () => {
  const { categories, loading } = useSelector((state) => state.categories);

  const [displayModal, setDisplayModal] = useState(false);
  const [selectedCategoryID, setSelectedCategoryID] = useState(null);

  const handleClose = () => setDisplayModal(false);
  const toggleModal = () => setDisplayModal(!displayModal);

  return (
    <Stack sx={{ py: 2 }}>
      <HeaderWithButton
        title="Categories"
        primaryButtonTextLabel="Add Category"
        primaryStartIcon={<AddRounded />}
        handleClickPrimaryButton={toggleModal}
      />
      <Category
        categories={categories}
        loading={loading}
        setSelectedCategoryID={setSelectedCategoryID}
        setDisplayModal={setDisplayModal}
      />
      {displayModal && (
        <SimpleModal title="Add New Category" handleClose={handleClose} maxSize="md">
          <AddCategory
            categories={categories}
            loading={loading}
            handleCloseAddCategory={handleClose}
            selectedCategoryID={selectedCategoryID}
            setSelectedCategoryID={setSelectedCategoryID}
          />
        </SimpleModal>
      )}
    </Stack>
  );
};

export default CategoryList;
