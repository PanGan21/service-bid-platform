import { useState, useEffect } from "react";
import styles from "./styles.module.css";

export const Pagination = ({
  pageChangeHandler,
  totalRows,
  rowsPerPage,
  currentPage,
}: {
  pageChangeHandler: Function;
  totalRows: number;
  rowsPerPage: number;
  currentPage: number;
}) => {
  // Calculating max number of pages
  const noOfPages = Math.ceil(totalRows / rowsPerPage);

  // Creating an array with length equal to no.of pages
  const pagesArr = [...new Array(noOfPages)];

  // Navigation arrows enable/disable state
  const [canGoBack, setCanGoBack] = useState(false);
  const [canGoNext, setCanGoNext] = useState(true);

  // These variables give the first and last record/row number
  // with respect to the current page
  const [pageFirstRecord, setPageFirstRecord] = useState(1);
  const [pageLastRecord, setPageLastRecord] = useState(rowsPerPage);

  // Onclick handlers for the butons
  const onNextPage = () => pageChangeHandler(currentPage + 1);
  const onPrevPage = () => pageChangeHandler(currentPage - 1);
  const onPageSelect = (pageNo: number) => pageChangeHandler(pageNo);

  // Disable previous and next buttons in the first and last page
  // respectively
  useEffect(() => {
    if (noOfPages === currentPage) {
      setCanGoNext(false);
    } else {
      setCanGoNext(true);
    }
    if (currentPage === 1) {
      setCanGoBack(false);
    } else {
      setCanGoBack(true);
    }
  }, [noOfPages, currentPage]);

  // To set the starting index of the page
  useEffect(() => {
    const skipFactor = (currentPage - 1) * rowsPerPage;
    // Some APIs require skip for paginaiton. If needed use that instead
    // pageChangeHandler(skipFactor);
    // pageChangeHandler(currentPage)
    setPageFirstRecord(skipFactor + 1);
  }, [currentPage, rowsPerPage]);

  // To set the last index of the page
  useEffect(() => {
    const count = pageFirstRecord + rowsPerPage;
    setPageLastRecord(count > totalRows ? totalRows : count - 1);
  }, [pageFirstRecord, rowsPerPage, totalRows]);

  return (
    <>
      {noOfPages > 1 ? (
        <div className={styles.pagination}>
          <div className={styles.pageInfo}>
            Showing {pageFirstRecord} - {pageLastRecord} of {totalRows}
          </div>
          <div className={styles.pagebuttons}>
            <button
              className={styles.pageBtn}
              onClick={onPrevPage}
              disabled={!canGoBack}
            >
              &#8249;
            </button>
            {pagesArr.map((num, index) => (
              <button
                key={index}
                onClick={() => onPageSelect(index + 1)}
                className={`${styles.pageBtn}  ${
                  index + 1 === currentPage ? styles.activeBtn : ""
                }`}
              >
                {index + 1}
              </button>
            ))}
            <button
              className={styles.pageBtn}
              onClick={onNextPage}
              disabled={!canGoNext}
            >
              &#8250;
            </button>
          </div>
        </div>
      ) : null}
    </>
  );
};
