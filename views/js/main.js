new Vue({
    el: '#app',
    delimiters: ['${', '}'],
    vuetify: new Vuetify(),
    data() {
        return {
            keys: [[1, 2, 3, 'C'], [4, 5, 6, 'AC'], [7, 8, 9]],
            table: this.newTable(null),
            editable: this.newTable(true),
            selectedRow: null,
            selectedCol: null,
            dialog: {
                requestAnswer: false,
                badTable: false,
                allClear: false
            },
            loading: false
        }
    },
    methods: {
        newTable: function(x) {
            var t = new Array(9)
            for (let i = 0; i < 9; i++) {
                t[i] = new Array(9).fill(x)
            }
            return t
        },
        selectCell: function(row, col) {
            this.selectedRow = row
            this.selectedCol = col
            return
        },
        isFocusedRow: function(row) {
            return row == this.selectedRow
        },
        isFocusedCol: function(col) {
            return col == this.selectedCol
        },
        isFocusedCell: function(row, col) {
            return this.isFocusedRow(row) && this.isFocusedCol(col)
        },
        setNumber: function(row, col, num) {
            if (0 <= row && row < 9 && 0 <= col && col < 9) {
                if (this.editable[row][col]) {
                    this.table[row].splice(col, 1, num)
                }
            }
        },
        allClear: function(strict) {
            this.selectCell(null, null)
            for (let i = 0; i < 9; i++) {
                for (let j = 0; j < 9; j++) {
                    if (!this.editable[i][j] && !strict) {
                        continue
                    }
                    this.table[i].splice(j, 1, null)
                    this.editable[i].splice(j, 1, true)
                }
            }
            this.dialog.allClear = false
        },
        keyAction: function(key) {
            let row = this.selectedRow
            let col = this.selectedCol
            if (key === 'C') {
                this.setNumber(row, col, null)
            } else if (key == 'AC') {
                this.dialog.allClear = true
            } else {
                this.setNumber(row, col, key)
            }
        },
        requestAnswer: function() {
            this.dialog.requestAnswer = false
            this.loading = true
            this.selectCell(null, null)
            // 0埋めを表示させないために新規テーブルを用意
            var sendTable = new Array(9)
            for (let i = 0; i < 9; i++) {
                sendTable[i] = new Array(9)
                for (let j = 0; j < 9; j++) {
                    sendTable[i][j] = this.table[i][j] || 0
                }
            }
            axios.post('/api/sudoku/solve', {
                table: sendTable
            }).then(res => {
                for (let i = 0; i < 9; i++) {
                    for (let j = 0; j < 9; j++) {
                        this.setNumber(i, j, res.data.answer[i][j])
                    }
                }
                this.loading = false
            }).catch(err => {
                console.log(err.response.data)
                this.loading = false
                this.dialog.badTable = true
            })
        }
    },
    mounted: function() {
        for (let i = 0; i < 9; i++) {
            this.table[i].splice(i, 1, i+1)
            this.editable[i].splice(i, 1, false)
        }
    }
})