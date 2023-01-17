import { useMemo } from "react";
import { useTable } from "react-table";
import { Loader } from "../loader/Loader";
import "./styles.module.css";

export type Column = {
  Header: string;
  accessor: string;
};

export const AppTable = ({
  columns,
  data,
  isLoading,
}: {
  columns: any;
  data: any;
  isLoading: boolean;
}) => {
  const columnData = useMemo(() => columns, [columns]);
  const rowData = useMemo(() => data, [data]);
  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } =
    useTable({
      columns: columnData,
      data: rowData,
    });
  return (
    <>
      {isLoading ? (
        <Loader />
      ) : (
        <>
          <table {...getTableProps()}>
            <thead>
              {headerGroups.map((headerGroup) => (
                <tr {...headerGroup.getHeaderGroupProps()}>
                  {headerGroup.headers.map((column) => (
                    <th {...column.getHeaderProps()}>
                      {column.render("Header")}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody {...getTableBodyProps()}>
              {rows.map((row, i) => {
                prepareRow(row);
                return (
                  <tr {...row.getRowProps()}>
                    {row.cells.map((cell) => {
                      return (
                        <td {...cell.getCellProps()}>{cell.render("Cell")}</td>
                      );
                    })}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </>
      )}
    </>
  );
};
