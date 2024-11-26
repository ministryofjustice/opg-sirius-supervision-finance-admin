package db

import (
	"context"
	"time"
)

const AgedDebtQuery = `WITH outstanding_invoices AS (SELECT i.id,
                                     i.finance_client_id,
                                     i.feetype,
                                     COALESCE(sl.supervision_level, i.feetype)                           AS supervision_level,
                                     i.reference,
                                     i.raiseddate,
                                     i.raiseddate + '30 days'::INTERVAL                                  AS due_date,
                                     (i.amount / 100.0)::NUMERIC(10, 2)                                  AS amount,
                                     ((i.amount - SUM(COALESCE(la.amount, 0))) / 100.00)::NUMERIC(10, 2) AS outstanding,
                                     EXTRACT(YEAR FROM
                                             AGE(NOW(), (i.raiseddate + '30 days'::INTERVAL)))::INT         age
                              FROM supervision_finance.invoice i
                                       LEFT JOIN supervision_finance.ledger_allocation la ON i.id = la.invoice_id
                                  AND la.status = 'ALLOCATED'
                                       LEFT JOIN LATERAL (
                                  SELECT ifr.supervisionlevel AS supervision_level
                                  FROM supervision_finance.invoice_fee_range ifr
                                  WHERE ifr.invoice_id = i.id
                                  ORDER BY id
                                  LIMIT 1
                                  ) sl ON TRUE
							WHERE i.raiseddate >= $1 AND i.raiseddate <= $2
                              GROUP BY i.id, i.amount, sl.supervision_level
                              HAVING i.amount > COALESCE(SUM(la.amount), 0)),
     age_per_client AS (SELECT fc.client_id, MAX(oi.age) AS age
                        FROM supervision_finance.finance_client fc
                                 JOIN outstanding_invoices oi ON fc.id = oi.finance_client_id
                        GROUP BY fc.client_id)
SELECT CONCAT(p.firstname, ' ', p.surname)                 AS "Customer Name",
       p.caserecnumber                                     AS "Customer number",
       fc.sop_number                                       AS "SOP number",
       d.deputytype                                        AS "Deputy type",
       COALESCE(active_orders.is_active, 'No')             AS "Active case?",
       '''0470'                                              AS "Entity",
       '99999999'                                          AS "Receivable cost centre",
       'BALANCE SHEET'                                     AS "Receivable cost centre description",
       '1816100000'                                        AS "Receivable account code",
       cc.code                                             AS "Revenue cost centre",
       cc.cost_centre_description                          AS "Revenue cost centre description",
       a.code                                              AS "Revenue account code",
       a.account_code_description                          AS "Revenue account code description",
       oi.feetype                                          AS "Invoice type",
       oi.reference                                        AS "Trx number",
       tt.description                                      AS "Transaction Description",
       oi.raiseddate                                       AS "Invoice date",
       oi.due_date                                         AS "Due date",
       CASE
           WHEN EXTRACT(MONTH FROM oi.raiseddate) >= 4 THEN EXTRACT(YEAR FROM oi.raiseddate)
           ELSE EXTRACT(YEAR FROM oi.raiseddate) - 1
           END                                             AS "Financial year",
       '30 NET'                                            AS "Payment terms",
       oi.amount                                           AS "Original amount",
       oi.outstanding                                      AS "Outstanding amount",
       CASE
           WHEN NOW() < (oi.due_date + '1 day'::INTERVAL) THEN oi.outstanding
           ELSE 0 END                                      AS "Current",
       CASE
           WHEN NOW() > oi.due_date AND oi.age < 2 THEN oi.outstanding
           ELSE 0 END                                      AS "0-1 years",
       CASE WHEN oi.age = 2 THEN oi.outstanding ELSE 0 END AS "1-2 years",
       CASE WHEN oi.age = 3 THEN oi.outstanding ELSE 0 END AS "2-3 years",
       CASE WHEN oi.age = 4 THEN oi.outstanding ELSE 0 END AS "3-5 years",
       CASE WHEN oi.age > 4 THEN oi.outstanding ELSE 0 END AS "5+ years",
       CASE
           WHEN apc.age < 2 THEN '''0-1'
           WHEN apc.age = 2 THEN '''1-2'
           WHEN apc.age = 3 THEN '''2-3'
           WHEN apc.age = 4 THEN '''3-5'
           ELSE '''5+' END                                   AS "Debt impairment years"
FROM supervision_finance.finance_client fc
         JOIN outstanding_invoices oi ON fc.id = oi.finance_client_id
         JOIN age_per_client apc ON fc.client_id = apc.client_id
         JOIN supervision_finance.transaction_type tt
              ON oi.feetype = tt.fee_type AND oi.supervision_level = tt.supervision_level
         JOIN supervision_finance.account a ON tt.account_code = a.code
         JOIN supervision_finance.cost_centre cc ON cc.code = a.cost_centre
         JOIN public.persons p ON fc.client_id = p.id
         LEFT JOIN public.persons d ON p.feepayer_id = d.id
         LEFT JOIN LATERAL (
    SELECT 'Yes' AS is_active
    FROM cases c
    WHERE p.id = c.client_id
      AND c.orderstatus = 'ACTIVE'
    LIMIT 1
    ) active_orders ON TRUE;`

func (c *Client) GetAgedDebt(ctx context.Context, fromDate time.Time, toDate time.Time) ([][]string, error) {
	items := [][]string{{
		"Customer Name",
		"Customer number",
		"SOP number",
		"Deputy type",
		"Active case?",
		"Entity",
		"Receivable cost centre",
		"Receivable cost centre description",
		"Receivable account code",
		"Revenue cost centre",
		"Revenue cost centre description",
		"Revenue account code",
		"Revenue account code description",
		"Invoice type",
		"Trx number",
		"Transaction Description",
		"Invoice date",
		"Due date",
		"Financial year",
		"Payment terms",
		"Original amount",
		"Outstanding amount",
		"Current",
		"0-1 years",
		"1-2 years",
		"2-3 years",
		"3-5 years",
		"5+ years",
		"Debt impairment years",
	}}

	rows, err := c.db.Query(ctx, AgedDebtQuery, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var i []string
		var stringValue string

		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			stringValue, _ = value.(string)
			i = append(i, stringValue)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
