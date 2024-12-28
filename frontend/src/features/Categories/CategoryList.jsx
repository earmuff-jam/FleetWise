import { useEffect, useState } from 'react';

import { useDispatch, useSelector } from 'react-redux';

import { Stack } from '@mui/material';

import SimpleModal from '@common/SimpleModal';
import AddCategory from '@features/Categories/AddCategory';
import SectionCardHeader from '@common/SectionCard/SectionCardHeader';
import { categoryActions } from '@features/Categories/categoriesSlice';
import SectionCardContent from '@common/SectionCard/SectionCardContent';

const CategoryList = ({ displayConcise = false }) => {
  const dispatch = useDispatch();
  const { categories = [], loading } = useSelector((state) => state.categories);

  const [sortedData, setSortedData] = useState([]);
  const [sortingOrder, setSortingOrder] = useState(true); // false ascending

  const [displayModal, setDisplayModal] = useState(false);
  const [selectedFilter, setSelectedFilter] = useState('');
  const [selectedCategoryID, setSelectedCategoryID] = useState('');

  const handleClose = () => {
    setDisplayModal(false);
    setSelectedCategoryID('');
  };

  const toggleModal = () => setDisplayModal(!displayModal);
  const handleDownload = () => dispatch(categoryActions.download());

  const removeSelectedCategory = (id) => dispatch(categoryActions.removeCategory({ id }));

  const filterAndBuildCategories = (displayConcise, categories, selectedFilter) => {
    if (displayConcise) {
      return categories?.slice(0, 4);
    } else if (selectedFilter.length > 0) {
      return categories.filter((element) => element.status_name === selectedFilter);
    } else {
      return sortedData;
    }
  };

  useEffect(() => {
    if (sortingOrder) {
      if (categories && categories.length > 0) {
        const draft = [...categories].sort((a, b) => new Date(a.updated_at) - new Date(b.updated_at));
        setSortedData(draft);
      } else {
        setSortedData(categories);
      }
    } else {
      setSortedData(categories);
    }
  }, [sortingOrder, categories]);

  useEffect(() => {
    dispatch(categoryActions.getCategories(100));
  }, []);

  return (
    <Stack sx={{ py: 2 }} data-tour="categories-0">
      <SectionCardHeader
        title="Categories"
        caption={selectedFilter ? `Applying ${selectedFilter} status filter` : 'Organize items into categories'}
        primaryBtnTitle="Add category"
        toggleModal={toggleModal}
        selectedFilter={selectedFilter}
        setSelectedFilter={setSelectedFilter}
        sortingOrder={sortingOrder}
        setSortingOrder={setSortingOrder}
        handleDownload={handleDownload}
        addBtnDataTour={'categories-1'}
        downloadBtnDataTour={'categories-2'}
        filterBtnDataTour={'categories-3'}
        sortBtnDataTour={'categories-4'}
        disableDownloadIcon={Boolean(categories) && categories.length <= 0}
      />
      <SectionCardContent
        loading={loading}
        displayModal={displayModal}
        setDisplayModal={setDisplayModal}
        setSelectedID={setSelectedCategoryID}
        removeItem={removeSelectedCategory}
        prefixURI={'category'}
        content={filterAndBuildCategories(displayConcise, categories, selectedFilter)}
      />
      {displayModal && (
        <SimpleModal
          title="Add new category"
          subtitle="Create categories to group assets. Assigned locations are approximate values."
          handleClose={handleClose}
          maxSize="sm"
        >
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
