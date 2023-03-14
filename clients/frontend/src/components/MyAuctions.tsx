import { useState, useEffect } from "react";
import { AppTable, Column } from "../common/table";
import { Pagination } from "../common/pagination";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countMyAuctions,
  formatAuctions,
  getMyAuctions,
} from "../services/auction";
import { FormattedAuction } from "../types/auction";
import { useNavigate } from "react-router-dom";

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
];

type Props = {};

export const MyAuctions: React.FC<Props> = () => {
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
  const navigate = useNavigate();

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countMyAuctions().then((response) => {
      if (response.data && response.data) {
        setTotalAuctions(response.data);
      }
    });

    getMyAuctions(ROWS_PER_TABLE_PAGE, currentPage).then((response) => {
      const auctions = response.data || [];
      setPageData({
        isLoading: false,
        rowData: formatAuctions(auctions),
        totalAuctions: totalAuctions,
      });
    });
  }, [currentPage, totalAuctions]);

  const handleRowSelection = (auction: any) => {
    const fullAuction = pageData.rowData.find((r) => r.Id === auction.Id);
    navigate("/assigned-auction", { state: fullAuction });
  };

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
