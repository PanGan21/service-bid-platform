import { useEffect, useState } from "react";
import { Pagination } from "../common/pagination";
import { AppTable, Column } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import { countMyBids, getMyBids } from "../services/bid";
import { Bid } from "../types/bid";

const columns: Column[] = [
  {
    Header: "Id",
    accessor: "Id",
  },
  {
    Header: "Amount €",
    accessor: "Amount",
  },
  {
    Header: "RequestId",
    accessor: "RequestId",
  },
];

type Props = {};

export const MyBids: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: Bid[];
    isLoading: boolean;
    totalBids: number;
  }>({
    rowData: [],
    isLoading: false,
    totalBids: 0,
  });
  const [totalBids, setTotalBids] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countMyBids().then((response) => {
      if (response.data && response.data) {
        setTotalBids(response.data);
      }
    });

    getMyBids(ROWS_PER_TABLE_PAGE, currentPage).then((response) => {
      const bids = response.data || [];
      setPageData({
        isLoading: false,
        rowData: bids,
        totalBids,
      });
    });
  }, [currentPage, totalBids]);

  const handleRowSelection = (request: any) => {};

  return (
    <div>
      <div style={{ height: "450px" }}>
        <AppTable
          columns={columns}
          data={pageData.rowData}
          isLoading={pageData.isLoading}
          onRowClick={(r) => handleRowSelection(r.values)}
        />
      </div>
      <Pagination
        totalRows={pageData.totalBids}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
