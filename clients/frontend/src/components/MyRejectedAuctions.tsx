import { useEffect, useState } from "react";
import { Column } from "react-table";
import { Pagination } from "../common/pagination";
import { AppTable } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countOwnRejectedAuctions,
  formatAuctions,
  getOwnRejectedAuctions,
} from "../services/auction";
import { FormattedAuction } from "../types/auction";

const columns: Column[] = [
  {
    Header: "Id",
    accessor: "Id",
  },
  {
    Header: "Title",
    accessor: "Title",
  },
  {
    Header: "Postcode",
    accessor: "Postcode",
  },
  {
    Header: "Info",
    accessor: "Info",
  },
  {
    Header: "Deadline",
    accessor: "Deadline",
  },
  {
    Header: "Status",
    accessor: "Status",
  },
  {
    Header: "Rejection reason",
    accessor: "RejectionReason",
  },
];

type Props = {};

export const MyRejectedAuctions: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: FormattedAuction[];
    isLoading: boolean;
    totalAuctions: number;
  }>({
    rowData: [],
    isLoading: false,
    totalAuctions: 0,
  });

  const [totalAuctions, setTotalAuctions] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countOwnRejectedAuctions().then((response) => {
      if (response.data && response.data) {
        setTotalAuctions(response.data);
      }
    });

    getOwnRejectedAuctions(ROWS_PER_TABLE_PAGE, currentPage).then(
      (response) => {
        const auctions = response.data || [];
        setPageData({
          isLoading: false,
          rowData: formatAuctions(auctions),
          totalAuctions: totalAuctions,
        });
      }
    );
  }, [currentPage, totalAuctions]);

  const handleRowSelection = (auction: any) => {};

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
        totalRows={pageData.totalAuctions}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
